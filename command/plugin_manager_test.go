package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/plugin"
)

const TestPluginProcessMagicCookie = "abcd"

func TestPluginLoad(t *testing.T) {
	plugin := testPlugin(t, "mock")
	if err := plugin.Load(); err != nil {
		t.Fatalf("err: %s", err)
	}

	if !reflect.DeepEqual(plugin.AppMeta, testPluginMockMeta) {
		t.Fatalf("bad: %#v", plugin.AppMeta)
	}
}

func TestPluginLoad_used(t *testing.T) {
	plugin := testPlugin(t, "mock")
	if err := plugin.Load(); err != nil {
		t.Fatalf("err: %s", err)
	}

	if plugin.Used() {
		t.Fatal("should not be used")
	}

	if _, err := plugin.App(); err != nil {
		t.Fatalf("err: %s", err)
	}

	if !plugin.Used() {
		t.Fatal("should be used")
	}
}

func TestPluginManager_saveLoad(t *testing.T) {
	mock := testPlugin(t, "mock")
	mgr := &PluginManager{
		plugins: []*Plugin{
			mock,
			testPlugin(t, "empty"),
			testPlugin(t, "empty"),
			testPlugin(t, "empty"),
		},
	}

	if err := mgr.LoadAll(); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Use only the mock
	if _, err := mock.App(); err != nil {
		t.Fatalf("err: %s", err)
	}
	if !mock.Used() {
		t.Fatal("should be used")
	}

	// Create the temporary file to save the data
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)
	path := filepath.Join(td, "used")

	if err := mgr.StoreUsed(path); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Create a new manager and load only the used
	mgr = &PluginManager{}
	if err := mgr.LoadUsed(path); err != nil {
		t.Fatalf("err: %s", err)
	}
	if v := mgr.Plugins(); len(v) != 1 {
		t.Fatalf("bad: %#v", v)
	}
}

// testPlugin returns a test plugin of the given name. This should correspond
// to one of the availabile plugins.
func testPlugin(t *testing.T, n string) *Plugin {
	var result Plugin
	result.Path = os.Args[0]
	result.Args = []string{
		"-test.run=TestPluginProcess",
		"--",
		TestPluginProcessMagicCookie,
		n,
	}

	return &result
}

func TestPluginProcess(*testing.T) {
	// Find where our arguments start based on "--"
	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}

		args = args[1:]
	}
	if len(args) < 1 {
		return
	}
	if args[0] != TestPluginProcessMagicCookie {
		return
	}

	// We know we're in a requested plugin process, we want to completely
	// exit at the end of this so let's defer that.
	defer os.Exit(0)

	// Determine what plugin we're supposed to be serving
	var opts plugin.ServeOpts
	switch args[1] {
	case "empty":
		appImpl := &app.Mock{}
		opts.AppFunc = func() app.App {
			return appImpl
		}
	case "mock":
		appImpl := &app.Mock{
			MetaResult: testPluginMockMeta,
		}
		opts.AppFunc = func() app.App {
			return appImpl
		}
	default:
		fmt.Fprintf(os.Stderr, "Invalid plugin: %s\n", args[1])
		os.Exit(2)
	}

	// Serve the plugin!
	plugin.Serve(&opts)
}

var testPluginMockMeta = &app.Meta{
	Tuples: []app.Tuple{
		app.Tuple{"test", "test", "test"},
	},
}
