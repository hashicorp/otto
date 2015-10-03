package localaddr

import (
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestDB_upgrade_v1(t *testing.T) {
	path := testCopyV1(t)
	defer os.Remove(path)

	// Get the CIDR that this should be in
	_, cidr, err := net.ParseCIDR("100.64.0.0/10")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Init the DB
	db := &DB{Path: path}

	// Grab an IP and verify it is in the proper place
	result, err := db.Next()
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if !cidr.Contains(result) {
		t.Fatal("is not in CIDR")
	}
}

func TestDB(t *testing.T) {
	td, err := ioutil.TempDir("", "localaddr")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

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

		results[result.String()] = struct{}{}
	}
}

func testCopyV1(t *testing.T) string {
	td, err := ioutil.TempDir("", "localaddr")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)
	path := filepath.Join(td, "addr.db")

	// Copy the v1 DB so we don't modify it directly
	src, err := os.Open("./test-fixtures/v1.db")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	dst, err := os.Create(path)
	if err != nil {
		src.Close()
		t.Fatalf("err: %s", err)
	}

	_, err = io.Copy(dst, src)
	src.Close()
	dst.Close()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return path
}
