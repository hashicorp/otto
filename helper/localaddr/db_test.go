package localaddr

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestDB(t *testing.T) {
	td, err := ioutil.TempDir("", "localaddr")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	// We put a random (but actually arbitrary) IP in here to test
	// Release later.
	var random net.IP

	// We shouldn't collide up to the max number
	results := make(map[string]struct{})
	db := &DB{Path: filepath.Join(td, "addr.db")}
	for i := 0; i < 254; i++ {
		result, err := db.Next()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if result.To4() == nil {
			t.Fatalf("not an IPv4: %s", result)
		}

		if _, ok := results[result.String()]; ok {
			t.Fatalf("collision: %s", result)
		}

		if i == 42 {
			random = result
		}

		results[result.String()] = struct{}{}
	}

	// The next one we get SHOULD collide
	actual, err := db.Next()
	if _, ok := results[actual.String()]; !ok {
		t.Fatal("should've collided")
	}

	// Release one and make sure we get that back
	if err := db.Release(random); err != nil {
		t.Fatalf("err: %s", err)
	}

	actual, err = db.Next()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if actual.String() != random.String() {
		t.Fatalf("bad: %s", actual)
	}
}
