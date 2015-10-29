package command

import (
	"fmt"
	"os"
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
