package appfile

import (
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/otto/appfile/detect"
)

// TestAppfile returns a compiled appfile for the given path. This uses
// defaults for detectors and such so it is up to you to use a fairly
// complete Appfile.
func TestAppfile(t TestT, path string) *Compiled {
	detector := &detect.Config{
		Detectors: []*detect.Detector{
			&detect.Detector{
				Type: "test",
				File: []string{"Appfile"},
			},
		},
	}

	// Create a loader that sets up defaults
	loader := func(f *File, dir string) (*File, error) {
		fDef, err := Default(dir, detector)
		if err != nil {
			return nil, err
		}

		// Default type should be "test"
		fDef.Infrastructure[0].Type = "test"
		fDef.Infrastructure[0].Flavor = "test"
		fDef.Infrastructure[0].Foundations = nil

		var merged File
		if err := merged.Merge(fDef); err != nil {
			return nil, err
		}
		if f != nil {
			if err := merged.Merge(f); err != nil {
				return nil, err
			}
		}

		return &merged, nil
	}

	// Parse the raw file
	f, err := ParseFile(path)
	if err != nil {
		t.Fatal("err: ", err)
	}

	// Call the loader for this Appfile
	f, err = loader(f, filepath.Dir(f.Path))
	if err != nil {
		t.Fatal("err: ", err)
	}

	// Create a temporary directory for the compilation data. We don't
	// delete this now in case we're using any of that data, but the
	// temp dir should get cleaned up by the system at some point.
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatal("err: ", err)
	}

	// Create the compiler
	compiler, err := NewCompiler(&CompileOpts{
		Dir:    td,
		Loader: loader,
	})
	if err != nil {
		t.Fatal("err: ", err)
	}

	// Compile
	result, err := compiler.Compile(f)
	if err != nil {
		t.Fatal("err: ", err)
	}

	return result
}

// TestT is the interface used to handle the test lifecycle of a test.
//
// Users should just use a *testing.T object, which implements this.
type TestT interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Skip(args ...interface{})
}
