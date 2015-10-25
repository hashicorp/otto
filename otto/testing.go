package otto

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/appfile/detect"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
)

// TestAppTuple is a basic app tuple that can be used for testing and
// has all fields set to "test".
var TestAppTuple = app.Tuple{"test", "test", "test"}

// TestAppfile returns a compiled appfile for the given path. This uses
// defaults for detectors and such so it is up to you to use a fairly
// complete Appfile.
func TestAppfile(t *testing.T, path string) *appfile.Compiled {
	def, err := appfile.Default(filepath.Dir(path), &detect.Config{
		Detectors: []*detect.Detector{
			&detect.Detector{
				Type: "test",
				File: []string{"Appfile"},
			},
		},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Default type should be "test"
	def.Infrastructure[0].Type = "test"
	def.Infrastructure[0].Flavor = "test"
	def.Infrastructure[0].Foundations = nil

	// Parse the raw file
	f, err := appfile.ParseFile(path)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Merge
	if err := def.Merge(f); err != nil {
		t.Fatalf("err: %s", err)
	}
	f = def

	// Create a temporary directory for the compilation data. We don't
	// delete this now in case we're using any of that data, but the
	// temp dir should get cleaned up by the system at some point.
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Compile it!
	result, err := appfile.Compile(f, &appfile.CompileOpts{
		Dir: td,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return result
}

// TestCoreConfig returns a CoreConfig that can be used for testing.
func TestCoreConfig(t *testing.T) *CoreConfig {
	// Temporary directory for data
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Basic config
	config := &CoreConfig{
		DataDir:    filepath.Join(td, "data"),
		LocalDir:   filepath.Join(td, "local"),
		CompileDir: filepath.Join(td, "compile"),
		Ui:         new(ui.Mock),
	}

	// Add some default mock implementations. These can be overwritten easily
	TestInfra(t, "test", config)
	TestApp(t, TestAppTuple, config)

	return config
}

// TestInfra adds a mock infrastructure with the given name to the
// core config and returns it.
func TestInfra(t *testing.T, n string, c *CoreConfig) *infrastructure.Mock {
	if c.Infrastructures == nil {
		c.Infrastructures = make(map[string]infrastructure.Factory)
	}

	result := new(infrastructure.Mock)
	c.Infrastructures[n] = func() (infrastructure.Infrastructure, error) {
		return result, nil
	}

	return result
}

// TestApp adds a mock app with the given tuple to the core config.
func TestApp(t *testing.T, tuple app.Tuple, c *CoreConfig) *app.Mock {
	if c.Apps == nil {
		c.Apps = make(map[app.Tuple]app.Factory)
	}

	result := new(app.Mock)
	c.Apps[tuple] = func() (app.App, error) {
		return result, nil
	}

	return result
}

// TestFoundation adds a mock foundation with the given tuple to the core config.
func TestFoundation(t *testing.T, tuple foundation.Tuple, c *CoreConfig) *foundation.Mock {
	if c.Foundations == nil {
		c.Foundations = make(map[foundation.Tuple]foundation.Factory)
	}

	result := new(foundation.Mock)
	c.Foundations[tuple] = func() (foundation.Foundation, error) {
		return result, nil
	}

	return result
}
