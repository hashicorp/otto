package localaddr

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestCachedDB(t *testing.T) {
	td, err := ioutil.TempDir("", "localaddr")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	db := &CachedDB{
		DB:        &DB{Path: filepath.Join(td, "addr.db")},
		CachePath: filepath.Join(td, "cache"),
	}

	ip, err := db.IP()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	next, err := db.IP()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if ip.String() != next.String() {
		t.Fatalf("bad: %s %s", next, ip)
	}
}
