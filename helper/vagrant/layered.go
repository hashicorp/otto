package vagrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/context"
)

// LayeredDir returns the directory where layered data is stored globally.
func LayeredDir(ctx *app.Context) string {
	return filepath.Join(ctx.GlobalDir, "vagrant-layered")
}

// Layered is a Vagrant environment that is created using a series of
// "layers". Otto manages these layers and this library automatically prunes
// unused layers. This library will also do the multi-process locking
// necessary to prevent races.
//
// To update a layer (change it), you should create a Layer with a new ID.
// IDs should be considered immutable for all time. This is to prevent breaking
// other environments. Once a layer is safely no longer in use by anybody
// for a sufficient period of time, Otto will automatically prune it.
//
// Layered itself doesn't manage the final Vagrant environment. This should
// be done outside of this using functions like Dev. Accounting should be done
// to avoid layers being pruned with `AddLeaf`, `RemoveLeaf`. If these
// aren't called layers underneath may be pruned which can corrupt leaves.
type Layered struct {
	// Layers are layers that are important for this run. This must include
	// all the Vagrantfiles for all the potential layers since we might need
	// to run all of them.
	Layers []*Layer

	// DataDir is the directory where Layered can write data to.
	DataDir string
}

// Layer is a single layer of the Layered Vagrant environment.
type Layer struct {
	// ID is a unique ID for the layer. See the note in Layered about
	// generating a new ID for every change/iteration in the Vagrantfile.
	ID string

	// Vagrantfile is the path to the Vagrantfile to bring up for this
	// layer. The Vagrantfile should handle all provisioning. This
	// Vagrantfile will be copied to another directory, so any paths
	// in it should be relative to the Vagrantfile.
	Vagrantfile string
}

// Build will build all the layers that are defined in this Layered
// struct. It will automatically output to the UI as needed.
//
// This will automatically acquire a process-lock to ensure that no duplicate
// layers are ever built. The process lock usually assumes that Otto is
// being run by the same user.
func (l *Layered) Build(ctx *context.Shared) error {
	// Grab the DB and initialize all the layers. This just inserts a
	// pending layer if it doesn't exist, as well as sets up the edges.
	db, err := l.db()
	if err != nil {
		return err
	}

	// We record the vertices so we can make use of them later
	layerVertices := make([]*layerVertex, len(l.Layers))
	for i, layer := range l.Layers {
		layerVertex, err := l.initLayer(db, layer)
		if err != nil {
			db.Close()
			return err
		}

		layerVertices[i] = layerVertex
		if i > 0 {
			// We have a prior layer, so setup the edge pointer
			err = db.Update(func(tx *bolt.Tx) error {
				bucket := tx.Bucket(boltEdgesBucket)
				return bucket.Put(
					[]byte(layer.ID),
					[]byte(layerVertices[i-1].Layer.ID))
			})
			if err != nil {
				return err
			}
		}
	}
	db.Close()

	// Go through each layer and build it. This will be a no-op if the
	// layer is already built.
	paths := l.LayerPaths()
	for _, layer := range l.Layers {
		path := paths[layer.ID]
		if err := l.buildLayer(layer, path, ctx); err != nil {
			return err
		}
	}

	return nil
}

func (l *Layered) buildLayer(layer *Layer, path string, ctx *context.Shared) error {
	// Layer isn't ready, so grab the lock on the layer and build it
	// TODO: multi-process lock

	// Once we have the lock, we check shortly in the DB if it is already
	// ready. If it is ready, we yield the lock and we're done!
	db, err := l.db()
	if err != nil {
		return err
	}
	layerV, err := l.readLayer(db, layer)
	db.Close()
	if err != nil {
		return err
	}
	if layerV.State == layerStateReady {
		return nil
	}

	// Prepare the build directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}

	// Copy the Vagrantfile into the destination path
	src, err := os.Open(layer.Vagrantfile)
	if err != nil {
		return err
	}
	dst, err := os.Create(filepath.Join(path, "Vagrantfile"))
	if err == nil {
		_, err = io.Copy(dst, src)
	}
	src.Close()
	dst.Close()
	if err != nil {
		return err
	}

	// Build the Vagrant instance
	vagrant := &Vagrant{
		Dir:     path,
		DataDir: path,
		Ui:      ctx.Ui,
	}
	if err := vagrant.Execute("up"); err != nil {
		return err
	}

	// Update the layer state that it is "ready"
	db, err = l.db()
	if err != nil {
		return err
	}
	defer db.Close()

	return l.updateLayer(db, layer, func(v *layerVertex) {
		v.State = layerStateReady
	})
}

// LayerPaths will return a mapping of all the paths that the Vagrantfiles
// should clone from. The key of the returned map is the ID of the layer
// and the value is the path.
//
// This can be used during the compilation process to setup proper paths.
func (l *Layered) LayerPaths() map[string]string {
	result := make(map[string]string)
	for _, layer := range l.Layers {
		result[layer.ID] = filepath.Join(l.DataDir, "layers", layer.ID)
	}

	return result
}

// db returns the database handle, and sets up the DB if it has never been created.
func (l *Layered) db() (*bolt.DB, error) {
	// Make the directory to store our DB
	if err := os.MkdirAll(l.DataDir, 0755); err != nil {
		return nil, err
	}

	// Create/Open the DB
	db, err := bolt.Open(filepath.Join(l.DataDir, "vagrant-layered.db"), 0644, nil)
	if err != nil {
		return nil, err
	}

	// Create the buckets
	err = db.Update(func(tx *bolt.Tx) error {
		for _, b := range boltBuckets {
			if _, err := tx.CreateBucketIfNotExists(b); err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Check the data version
	var version byte
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltVagrantBucket)
		data := bucket.Get([]byte("version"))
		if data == nil || len(data) == 0 {
			version = boltDataVersion
			return bucket.Put([]byte("version"), []byte{boltDataVersion})
		}

		version = data[0]
		return nil
	})
	if err != nil {
		return nil, err
	}

	if version > boltDataVersion {
		return nil, fmt.Errorf(
			"Vagrant layer data version is higher than this version of Otto knows how\n"+
				"to handle! This version of Otto can read up to version %d,\n"+
				"but version %d data file found.\n\n"+
				"This means that a newer version of Otto touched this data,\n"+
				"or the data was corrupted in some other way.",
			boltDataVersion, version)
	}

	return db, nil
}

// initLayer sets up the layer in the database
func (l *Layered) initLayer(db *bolt.DB, layer *Layer) (*layerVertex, error) {
	var result layerVertex
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLayersBucket)
		key := []byte(layer.ID)
		data := bucket.Get(key)
		if len(data) > 0 {
			return l.structRead(&result, data)
		}

		// Vertex doesn't exist. Create it and save it
		result = layerVertex{Layer: layer, State: layerStatePending}
		data, err := l.structData(&result)
		if err != nil {
			return err
		}

		// Write the pending layer
		return bucket.Put(key, data)
	})

	return &result, err
}

func (l *Layered) readLayer(db *bolt.DB, layer *Layer) (*layerVertex, error) {
	var result layerVertex
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLayersBucket)
		key := []byte(layer.ID)
		data := bucket.Get(key)
		if len(data) > 0 {
			return l.structRead(&result, data)
		}

		return fmt.Errorf("layer %s not found", layer.ID)
	})

	return &result, err
}

func (l *Layered) updateLayer(db *bolt.DB, layer *Layer, f func(*layerVertex)) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLayersBucket)
		key := []byte(layer.ID)
		data := bucket.Get(key)
		if len(data) == 0 {
			// This should never happen through this struct
			panic(fmt.Errorf("layer %s not found", layer.ID))
		}

		// Read the vertex, call the function to modify it
		var v layerVertex
		if err := l.structRead(&v, data); err != nil {
			return err
		}
		f(&v)

		// Save the resulting layer data
		data, err := l.structData(&v)
		if err != nil {
			return err
		}
		return bucket.Put(key, data)
	})
}

func (l *Layered) structData(d interface{}) ([]byte, error) {
	// Let's just output it in human-readable format to make it easy
	// for debugging. Disk space won't matter that much for this data.
	return json.MarshalIndent(d, "", "\t")
}

func (l *Layered) structRead(d interface{}, raw []byte) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	return dec.Decode(d)
}

var (
	boltVagrantBucket = []byte("vagrant")
	boltLayersBucket  = []byte("layers")
	boltEdgesBucket   = []byte("edges")
	boltEnvsBucket    = []byte("envs")
	boltBuckets       = [][]byte{
		boltVagrantBucket,
		boltLayersBucket,
		boltEdgesBucket,
		boltEnvsBucket,
	}
)

var (
	boltDataVersion byte = 1
)

// layerVertex is the type of vertex in the graph that is used to track
// layer usage throughout Otto.
type layerVertex struct {
	Layer *Layer     `json:"layer"`
	State layerState `json:"state"`
}

type layerState byte

const (
	layerStateInvalid layerState = iota
	layerStatePending
	layerStateReady
)
