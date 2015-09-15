package directory

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func TestBoltBackend(t *testing.T) {
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	TestBackend(t, &BoltBackend{
		Dir: td,
	})
}

func TestBoltBackend_version(t *testing.T) {
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	// Grab the DB and manually insert a version (breaking a black box here)
	b := &BoltBackend{Dir: td}
	db, err := b.db()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(boltOttoBucket)
		return b.Put([]byte("version"), []byte{boltDataVersion + 1})
	})
	db.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Re-instantiate and try an operation, it should fail
	b = &BoltBackend{Dir: td}
	db, err = b.db()
	if err == nil {
		db.Close()
		t.Fatal("should error")
	}
}
