package scriptpack

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	testBasic "github.com/hashicorp/otto/scriptpack/test-basic"
)

func TestScriptPackWrite(t *testing.T) {
	// Build the SP
	sp := &ScriptPack{
		Name: "foo",
		Data: testBasic.Bindata,
	}

	// Temporary dir to test contents
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	// Write the data!
	if err := sp.Write(td); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify some files exist
	var path string
	path = filepath.Join(td, sp.Name, "hello.txt")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("err: %s", err)
	}
}
