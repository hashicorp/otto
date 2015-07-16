package bindata

import (
	"io/ioutil"
	"os"
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
}

func testData() *Data {
	return &Data{
		Asset:    Asset,
		AssetDir: AssetDir,
	}
}
