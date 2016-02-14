package context

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/ui"
)

// TestShared returns an instance of Shared for testing with mock data.
func TestShared(t *testing.T) *Shared {
	tempDir, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	directory := &directory.BoltBackend{Dir: filepath.Join(tempDir, "directory")}

	return &Shared{
		Ui:         &ui.Mock{},
		Directory:  directory,
		InstallDir: filepath.Join(tempDir, "install-dir"),
	}
}
