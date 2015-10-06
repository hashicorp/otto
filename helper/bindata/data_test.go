package bindata

import (
	"bytes"
	"crypto/md5"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDataCopyDir(t *testing.T) {
	d := testData()

	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	if err := d.CopyDir(td, "test-data/copy-dir-basic"); err != nil {
		t.Fatalf("err: %s", err)
	}

	testCompareDir(t, td, "test-data/copy-dir-basic-expected")
}

func TestDataCopyDir_extends(t *testing.T) {
	d := testData()

	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	if err := d.CopyDir(td, "test-data/copy-dir-extends"); err != nil {
		t.Fatalf("err: %s", err)
	}

	testCompareDir(t, td, "test-data/copy-dir-extends-expected")
}

func testData() *Data {
	return &Data{
		Asset:    Asset,
		AssetDir: AssetDir,
		Context: map[string]interface{}{
			"value": "foo",
		},
	}
}

func testCompareDir(t *testing.T, a string, b string) {
	f, err := os.Open(a)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer f.Close()

	if _, err := os.Stat(b); err != nil {
		t.Fatalf("doesn't exist: %s", err)
	}

	fis, err := f.Readdir(-1)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	for _, fi := range fis {
		aPath := filepath.Join(a, fi.Name())
		bPath := filepath.Join(b, fi.Name())
		if filepath.Ext(bPath) == ".tpl" {
			bPath = bPath[:len(bPath)-len(".tpl")]
		}

		if fi.IsDir() {
			testCompareDir(t, aPath, bPath)
			continue
		}

		hash := md5.New()

		// Hash A
		var aReal bytes.Buffer
		af, err := os.Open(aPath)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		_, err = io.Copy(io.MultiWriter(&aReal, hash), af)
		af.Close()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		aHash := hash.Sum(nil)

		// Hash B
		hash.Reset()
		bf, err := os.Open(bPath)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		_, err = io.Copy(hash, bf)
		bf.Close()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		bHash := hash.Sum(nil)

		if !reflect.DeepEqual(aHash, bHash) {
			t.Fatalf("difference: %s\n\n%s", aPath, aReal.String())
		}
	}
}
