package directory

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/boltdb/bolt"
)

var (
	boltOttoBucket  = []byte("otto")
	boltAppsBucket  = []byte("apps")
	boltBlobBucket  = []byte("blob")
	boltInfraBucket = []byte("infra")
	boltBuckets     = [][]byte{
		boltOttoBucket,
		boltAppsBucket,
		boltBlobBucket,
		boltInfraBucket,
	}
)

var (
	boltDataVersion byte = 1
)

// BoltBackend is a Directory backend that stores data on local disk
// using BoltDB.
//
// The primary use case for the BoltBackend is out-of-box experience
// for Otto and single developers. For team usage, BoltBackend is not
// recommended.
//
// This backend also implements io.Closer and should be closed.
type BoltBackend struct {
	// Directory where data will be written. This directory will be
	// created if it doesn't exist.
	Dir string

	raw     *bolt.DB
	rawLock sync.Mutex
}

func (b *BoltBackend) Close() error {
	b.rawLock.Lock()
	defer b.rawLock.Unlock()

	if b.raw != nil {
		err := b.raw.Close()
		b.raw = nil
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *BoltBackend) GetBlob(k string) (*BlobData, error) {
	db, err := b.db()
	if err != nil {
		return nil, err
	}

	var data []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltBlobBucket)
		data = bucket.Get([]byte(k))
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &BlobData{
		Key:  k,
		Data: bytes.NewReader(data),
	}, nil
}

func (b *BoltBackend) PutBlob(k string, d *BlobData) error {
	db, err := b.db()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, d.Data); err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltBlobBucket)
		return bucket.Put([]byte(k), buf.Bytes())
	})
}

func (b *BoltBackend) GetInfra(infra *Infra) (*Infra, error) {
	db, err := b.db()
	if err != nil {
		return nil, err
	}

	var result *Infra
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltInfraBucket).Bucket([]byte(
			infra.Lookup.Infra))

		// If the bucket doesn't exist, we haven't written this yet
		if bucket == nil {
			return nil
		}

		// Get the key for this infra
		data := bucket.Get([]byte(b.infraKey(infra)))
		if data == nil {
			return nil
		}

		result = &Infra{}
		return b.structRead(result, data)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *BoltBackend) PutInfra(infra *Infra) error {
	if infra.ID == "" {
		infra.setId()
	}

	db, err := b.db()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		data, err := b.structData(infra)
		if err != nil {
			return err
		}

		bucket := tx.Bucket(boltInfraBucket)
		bucket, err = bucket.CreateBucketIfNotExists([]byte(
			infra.Lookup.Infra))
		if err != nil {
			return err
		}

		return bucket.Put([]byte(b.infraKey(infra)), data)
	})
}

func (b *BoltBackend) GetBuild(build *Build) (*Build, error) {
	db, err := b.db()
	if err != nil {
		return nil, err
	}

	var result *Build
	err = db.View(func(tx *bolt.Tx) error {
		// Get the app bucket
		bucket := tx.Bucket(boltAppsBucket).Bucket([]byte(
			build.Lookup.AppID))
		if bucket == nil {
			return nil
		}

		// Get the infra bucket
		bucket = bucket.Bucket([]byte(fmt.Sprintf(
			"%s-%s", build.Lookup.Infra, build.Lookup.InfraFlavor)))
		if bucket == nil {
			return nil
		}

		// Get the key for this infra
		data := bucket.Get([]byte("build"))
		if data == nil {
			return nil
		}

		result = &Build{}
		return b.structRead(result, data)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *BoltBackend) PutBuild(build *Build) error {
	db, err := b.db()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		data, err := b.structData(build)
		if err != nil {
			return err
		}

		// Get the app bucket
		bucket := tx.Bucket(boltAppsBucket)
		bucket, err = bucket.CreateBucketIfNotExists([]byte(
			build.Lookup.AppID))
		if err != nil {
			return err
		}

		// Get the infra bucket
		bucket, err = bucket.CreateBucketIfNotExists([]byte(fmt.Sprintf(
			"%s-%s", build.Lookup.Infra, build.Lookup.InfraFlavor)))
		if err != nil {
			return err
		}

		return bucket.Put([]byte("build"), data)
	})
}

func (b *BoltBackend) GetDeploy(deploy *Deploy) (*Deploy, error) {
	db, err := b.db()
	if err != nil {
		return nil, err
	}

	var result *Deploy
	err = db.View(func(tx *bolt.Tx) error {
		// Get the app bucket
		bucket := tx.Bucket(boltAppsBucket).Bucket([]byte(
			deploy.Lookup.AppID))
		if bucket == nil {
			return nil
		}

		// Get the infra bucket
		bucket = bucket.Bucket([]byte(fmt.Sprintf(
			"%s-%s", deploy.Lookup.Infra, deploy.Lookup.InfraFlavor)))
		if bucket == nil {
			return nil
		}

		// Get the key for this infra
		data := bucket.Get([]byte("deploy"))
		if data == nil {
			return nil
		}

		result = &Deploy{}
		return b.structRead(result, data)
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *BoltBackend) PutDeploy(deploy *Deploy) error {
	if deploy.ID == "" {
		deploy.setId()
	}

	db, err := b.db()
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		data, err := b.structData(deploy)
		if err != nil {
			return err
		}

		// Get the app bucket
		bucket := tx.Bucket(boltAppsBucket)
		bucket, err = bucket.CreateBucketIfNotExists([]byte(
			deploy.Lookup.AppID))
		if err != nil {
			return err
		}

		// Get the infra bucket
		bucket, err = bucket.CreateBucketIfNotExists([]byte(fmt.Sprintf(
			"%s-%s", deploy.Lookup.Infra, deploy.Lookup.InfraFlavor)))
		if err != nil {
			return err
		}

		return bucket.Put([]byte("deploy"), data)
	})
}

func (b *BoltBackend) infraKey(infra *Infra) string {
	key := "root"
	if infra.Lookup.Foundation != "" {
		key = fmt.Sprintf("foundation-%s", infra.Lookup.Foundation)
	}

	return key
}

// db returns the database handle, and sets up the DB if it has never
// been created.
func (b *BoltBackend) db() (*bolt.DB, error) {
	b.rawLock.Lock()
	defer b.rawLock.Unlock()

	if b.raw != nil {
		return b.raw, nil
	}

	// Make the directory to store our DB
	if err := os.MkdirAll(b.Dir, 0755); err != nil {
		return nil, err
	}

	// Create/Open the DB
	var err error
	b.raw, err = bolt.Open(filepath.Join(b.Dir, "otto.db"), 0644, nil)
	if err != nil {
		return nil, err
	}

	// Create the buckets
	err = b.raw.Update(func(tx *bolt.Tx) error {
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

	// Check the Otto version
	var version byte
	err = b.raw.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(boltOttoBucket)
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
			"Data version is higher than this version of Otto knows how\n"+
				"to handle! This version of Otto can read up to version %d,\n"+
				"but version %d data file found.\n\n"+
				"This means that a newer version of Otto touched this data,\n"+
				"or the data was corrupted in some other way.",
			boltDataVersion, version)
	}

	return b.raw, nil
}

func (b *BoltBackend) structData(d interface{}) ([]byte, error) {
	// Let's just output it in human-readable format to make it easy
	// for debugging. Disk space won't matter that much for this data.
	return json.MarshalIndent(d, "", "\t")
}

func (b *BoltBackend) structRead(d interface{}, raw []byte) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	return dec.Decode(d)
}
