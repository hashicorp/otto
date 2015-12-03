package otto

import (
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/appfile/detect"
)

// TestAppfile returns a compiled appfile for the given path. This uses
// defaults for detectors and such so it is up to you to use a fairly
// complete Appfile.
func TestAppfile(t TestT, path string) *appfile.Compiled {
	def, err := appfile.Default(filepath.Dir(path), &detect.Config{
		Detectors: []*detect.Detector{
			&detect.Detector{
				Type: "test",
				File: []string{"Appfile"},
			},
		},
	})
	if err != nil {
		t.Fatal("err: ", err)
	}

	// Default type should be "test"
	def.Infrastructure[0].Type = "test"
	def.Infrastructure[0].Flavor = "test"
	def.Infrastructure[0].Foundations = nil

	// Parse the raw file
	f, err := appfile.ParseFile(path)
	if err != nil {
		t.Fatal("err: ", err)
	}

	// Merge
	if err := def.Merge(f); err != nil {
		t.Fatal("err: ", err)
	}
	f = def

	// Create a temporary directory for the compilation data. We don't
	// delete this now in case we're using any of that data, but the
	// temp dir should get cleaned up by the system at some point.
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatal("err: ", err)
	}

	// Compile it!
	compiler, err := appfile.NewCompiler(&appfile.CompileOpts{
		Dir: td,
	})
	if err != nil {
		t.Fatal("err: ", err)
	}
	result, err := compiler.Compile(f)
	if err != nil {
		t.Fatal("err: ", err)
	}

	return result
}
