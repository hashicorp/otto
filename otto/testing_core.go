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
