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

	"github.com/hashicorp/otto/helper/bindata/test-pkg"
)

func TestDataCopyDir(t *testing.T) {
	cases := []struct {
		Dir string
	}{
		{"copy-dir-basic"},
		{"copy-dir-extends"},
		{"copy-dir-extends-var"},
		{"copy-dir-include"},
		{"copy-dir-shared"},
	}

	for _, tc := range cases {
		d := testData()

		func() {
			td, err := ioutil.TempDir("", "otto")
			if err != nil {
				t.Fatalf("err: %s", err)
			}
			defer os.RemoveAll(td)

			dir := filepath.Join("test-data", tc.Dir)
			if err := d.CopyDir(td, dir); err != nil {
				t.Fatalf("err: %s", err)
			}

			testCompareDir(t, td, dir+"-expected")
		}()
	}
}

func testData() *Data {
	return &Data{
		Asset:    Asset,
		AssetDir: AssetDir,
		Context: map[string]interface{}{
			"value": "foo",
		},
		SharedExtends: map[string]*Data{
			"foo": &Data{
				Asset:    testpkg.Asset,
				AssetDir: testpkg.AssetDir,
			},
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
