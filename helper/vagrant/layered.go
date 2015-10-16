package vagrant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/terraform/dag"
)

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
	vs, err := l.init(db)
	db.Close()
	if err != nil {
		return err
	}

	// Go through each layer and build it. This will be a no-op if the
	// layer is already built.
	for i, v := range vs {
		var last *layerVertex
		if i > 0 {
			last = vs[i-1]
		}

		if err := l.buildLayer(v, last, ctx); err != nil {
			return err
		}
	}

	return nil
}

// Prune will destroy all layers that haven't been used in a certain
// amount of time.
//
// TODO: "certain amount of time" for now we just prune any orphans
func (l *Layered) Prune(ctx *context.Shared) error {
	db, err := l.db()
	if err != nil {
		return err
	}
	defer db.Close()

	graph, err := l.graph(db)
	if err != nil {
		return err
	}

	// Get all the bad roots. These are anything without something depending
	// on it except for the main "root"
	roots := make([]dag.Vertex, 0)
	for _, v := range graph.Vertices() {
		if v == "root" {
			continue
		}
		if graph.UpEdges(v).Len() == 0 {
			roots = append(roots, v)
		}
	}
	if len(roots) == 0 {
		return nil
	}

	// Go through the remaining roots, these are the environments
	// that must be destroyed.
	for _, root := range roots {
		if err := l.pruneLayer(db, root.(*layerVertex), ctx); err != nil {
			return err
		}
	}

	return nil
}

// AddEnv will store the given environment as a user of this layer set,
// preventing the pruning of the layers here.
//
// This will also modify the argument to set the environment variable
// to point to the proper layer.
func (l *Layered) AddEnv(v *Vagrant) error {
	// Get the final layer
	layer := l.Layers[len(l.Layers)-1]

	// Update the DB with our environment
	db, err := l.db()
	if err != nil {
		return err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltEnvsBucket)
		key := []byte(v.DataDir)
		return bucket.Put(key, []byte(layer.ID))
	})
	db.Close()
	if err != nil {
		return err
	}

	// Get the path for the final layer and add it to the environment
	path := filepath.Join(l.layerPath(layer), "Vagrantfile")
	if v.Env == nil {
		v.Env = make(map[string]string)
	}
	v.Env[layerPathEnv] = path

	return nil
}

// RemoveEnv will remove the environment from the tracked layers.
func (l *Layered) RemoveEnv(v *Vagrant) error {
	db, err := l.db()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltEnvsBucket)
		key := []byte(v.DataDir)
		return bucket.Delete(key)
	})
}

// Pending returns a list of layers that are pending creation.
// Note that between calling this and calling something like Build(),
// this state may be different.
func (l *Layered) Pending() ([]string, error) {
	// Grab the DB and initialize all the layers. This just inserts a
	// pending layer if it doesn't exist, as well as sets up the edges.
	db, err := l.db()
	if err != nil {
		return nil, err
	}
	vs, err := l.init(db)
	db.Close()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0, len(vs))
	for _, v := range vs {
		if v.State != layerStateReady {
			result = append(result, v.Layer.ID)
		}
	}

	return result, nil
}

func (l *Layered) buildLayer(v *layerVertex, lastV *layerVertex, ctx *context.Shared) error {
	layer := v.Layer
	path := v.Path

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

	// Tell the user things are happening
	ctx.Ui.Header(fmt.Sprintf("Creating layer: %s", layer.ID))

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

	// Build the Vagrant instance. We bring it up, and then immediately
	// shut it down since we don't need it running. We start by trying to
	// destroy it in case there is another prior instance here.
	vagrant := &Vagrant{
		Dir:     path,
		DataDir: filepath.Join(path, ".vagrant"),
		Ui:      ctx.Ui,
	}
	if lastV != nil {
		vagrant.Env = map[string]string{
			layerPathEnv: filepath.Join(lastV.Path, "Vagrantfile"),
		}
	}
	if err := vagrant.Execute("destroy", "-f"); err != nil {
		return err
	}
	if err := vagrant.Execute("up"); err != nil {
		return err
	}
	if err := vagrant.Execute("halt"); err != nil {
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

func (l *Layered) pruneLayer(db *bolt.DB, v *layerVertex, ctx *context.Shared) error {
	layer := v.Layer
	path := v.Path

	ctx.Ui.Header(fmt.Sprintf(
		"Deleting layer '%s'...", layer.ID))

	// First, note that the layer is no longer ready
	err := l.updateLayer(db, layer, func(v *layerVertex) {
		v.State = layerStatePending
	})
	if err != nil {
		return err
	}

	// Check the path. If the path doesn't exist, then it is already destroyed.
	// If the path does exist, then we do an actual vagrant destroy
	_, err = os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		vagrant := &Vagrant{
			Dir:     path,
			DataDir: filepath.Join(path, ".vagrant"),
			Ui:      ctx.Ui,
		}
		if err := vagrant.Execute("destroy", "-f"); err != nil {
			return err
		}
	}

	// Delete the layer
	return l.deleteLayer(db, layer)
}

func (l *Layered) layerPath(layer *Layer) string {
	return filepath.Join(l.DataDir, "layers", layer.ID)
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

// init initializes the database for this layer setup.
func (l *Layered) init(db *bolt.DB) ([]*layerVertex, error) {
	layerVertices := make([]*layerVertex, len(l.Layers))
	for i, layer := range l.Layers {
		var parent *Layer
		if i > 0 {
			parent = l.Layers[i-1]
		}

		layerVertex, err := l.initLayer(db, layer, parent)
		if err != nil {
			return nil, err
		}

		layerVertices[i] = layerVertex
		if parent != nil {
			// We have a prior layer, so setup the edge pointer
			err = db.Update(func(tx *bolt.Tx) error {
				bucket := tx.Bucket(boltEdgesBucket)
				return bucket.Put(
					[]byte(layer.ID),
					[]byte(parent.ID))
			})
			if err != nil {
				return nil, err
			}
		}
	}

	return layerVertices, nil
}

// initLayer sets up the layer in the database
func (l *Layered) initLayer(db *bolt.DB, layer *Layer, parent *Layer) (*layerVertex, error) {
	var parentID string
	if parent != nil {
		parentID = parent.ID
	}

	var result layerVertex
	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLayersBucket)
		key := []byte(layer.ID)
		data := bucket.Get(key)
		if len(data) > 0 {
			var v layerVertex
			if err := l.structRead(&v, data); err != nil {
				return err
			}

			if v.Parent == parentID {
				result = v
				return nil
			}

			// The parent didn't match, so we just initialize a new
			// entry below. This will also force the destruction of the
			// old environment.
		}

		// Vertex doesn't exist. Create it and save it
		result = layerVertex{
			Layer:  layer,
			State:  layerStatePending,
			Parent: parent.ID,
			Path:   l.layerPath(layer),
		}
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

func (l *Layered) deleteLayer(db *bolt.DB, layer *Layer) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLayersBucket)
		key := []byte(layer.ID)
		return bucket.Delete(key)
	})
}

func (l *Layered) graph(db *bolt.DB) (*dag.AcyclicGraph, error) {
	graph := new(dag.AcyclicGraph)
	graph.Add("root")

	// First, add all the layers
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltLayersBucket)
		return bucket.ForEach(func(k, data []byte) error {
			var v layerVertex
			if err := l.structRead(&v, data); err != nil {
				return err
			}

			graph.Add(&v)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	// Next, connect the layers
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltEdgesBucket)
		return bucket.ForEach(func(k, data []byte) error {
			from := &layerVertex{Layer: &Layer{ID: string(k)}}
			to := &layerVertex{Layer: &Layer{ID: string(data)}}
			graph.Connect(dag.BasicEdge(from, to))
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	// Finally, add and connect all the envs
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltEnvsBucket)
		return bucket.ForEach(func(k, data []byte) error {
			key := fmt.Sprintf("env-%s", string(k))
			graph.Add(key)

			to := &layerVertex{Layer: &Layer{ID: string(data)}}
			graph.Connect(dag.BasicEdge(key, to))
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	return graph, nil
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

const layerPathEnv = "OTTO_VAGRANT_LAYER_PATH"

// layerVertex is the type of vertex in the graph that is used to track
// layer usage throughout Otto.
type layerVertex struct {
	Layer  *Layer     `json:"layer"`
	State  layerState `json:"state"`
	Parent string     `json:"parent"`
	Path   string     `json:"path"`
}

func (v *layerVertex) Hashcode() string {
	return fmt.Sprintf("layer-%s", v.Layer.ID)
}

func (v *layerVertex) Name() string {
	return v.Layer.ID
}

type layerState byte

const (
	layerStateInvalid layerState = iota
	layerStatePending
	layerStateReady
)
