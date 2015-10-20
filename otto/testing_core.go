package otto

import (
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/ui"
)

// TestAppTuple is a basic app tuple that can be used for testing and
// has all fields set to "test".
var TestAppTuple = app.Tuple{"test", "test", "test"}

// TestCoreOpts is a specialized struct that is used to create a Core,
// focused on the most common usage for tests.
type TestCoreOpts struct {
	// Path is the path to an Appfile to compile
	Path string

	// App to register with the TestAppTuple as a fixed result
	App app.App
}

// TestCore returns a *Core for testing. If TestCoreOpts is nil then
// this is equivalent to creating a core with TestCoreConfig set.
func TestCore(t TestT, config *TestCoreOpts) *Core {
	// Get the base config because we'll need this anyways
	coreConfig := TestCoreConfig(t)

	// If a config is set, then use that to do things
	if config != nil {
		if config.Path != "" {
			coreConfig.Appfile = TestAppfile(t, config.Path)
		}

		if config.App != nil {
			coreConfig.Apps[TestAppTuple] = func() (app.App, error) {
				return config.App, nil
			}
		}
	}

	// Create the core!
	core, err := NewCore(coreConfig)
	if err != nil {
		t.Fatal("error creating core: ", err)
	}

	return core
}

// TestCoreConfig returns a CoreConfig that can be used for testing.
func TestCoreConfig(t TestT) *CoreConfig {
	// Temporary directory for data
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatal("err: ", err)
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
