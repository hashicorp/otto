package bindata

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

	if _, err := ioutil.ReadFile(filepath.Join(td, "t")); err != nil {
		t.Fatalf("expected t.tpl to be rendered as t, err: %s", err)
	}
}

func testData() *Data {
	return &Data{
		Asset:    Asset,
		AssetDir: AssetDir,
	}
}
