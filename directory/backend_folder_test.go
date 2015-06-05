package directory

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFolderBackend(t *testing.T) {
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	TestBackend(t, &FolderBackend{
		Dir: td,
	})
}
