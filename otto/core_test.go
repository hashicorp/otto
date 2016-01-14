package otto

import (
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/directory"
)

func TestCoreApp(t *testing.T) {
	// Make a core that returns a fixed app
	coreConfig := TestCoreConfig(t)
	coreConfig.Appfile = TestAppfile(t, testPath("basic", "Appfile"))
	appMock := TestApp(t, TestAppTuple, coreConfig)
	core := testCore(t, coreConfig)

	// Get the App
	app, _, err := core.App()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if app != appMock {
		t.Fatal("should be equal")
	}
}

func TestCoreCompile_close(t *testing.T) {
	// Make a core that returns a fixed app
	coreConfig := TestCoreConfig(t)
	coreConfig.Appfile = TestAppfile(t, testPath("basic", "Appfile"))
	appMock := TestApp(t, TestAppTuple, coreConfig)
	core := testCore(t, coreConfig)

	// Compile!
	if err := core.Compile(); err != nil {
		t.Fatalf("err: %s", err)
	}

	if !appMock.CompileCalled {
		t.Fatal("compile should be called")
	}
	if !appMock.CloseCalled {
		t.Fatal("close should be called")
	}
}

func TestCoreCompile_customizationFilter(t *testing.T) {
	// Make a core that returns a fixed app
	coreConfig := TestCoreConfig(t)
	coreConfig.Appfile = TestAppfile(t, testPath("customization-app-filter", "Appfile"))
	appMock := TestApp(t, TestAppTuple, coreConfig)
	core := testCore(t, coreConfig)

	// Compile!
	if err := core.Compile(); err != nil {
		t.Fatalf("err: %s", err)
	}

	if !appMock.CompileCalled {
		t.Fatal("compile should be called")
	}

	// Verify our customizations
	var keys []string
	for _, c := range appMock.CompileContext.Appfile.Customization.Raw {
		for k, _ := range c.Config {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	expected := []string{"bar", "foo"}
	if !reflect.DeepEqual(keys, expected) {
		t.Fatalf("bad: %#v", keys)
	}
}

func TestCoreCompile_directory(t *testing.T) {
	// Make a core that returns a fixed app
	coreConfig := TestCoreConfig(t)
	coreConfig.Appfile = TestAppfile(t, testPath("compile-directory", "Appfile"))
	appMock := TestApp(t, TestAppTuple, coreConfig)
	core := testCore(t, coreConfig)

	// Compile!
	if err := core.Compile(); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify basic compile behavior
	if !appMock.CompileCalled {
		t.Fatal("compile should be called")
	}

	// Verify we're now in the directory
	d := coreConfig.Directory
	app, err := d.GetApp(directory.AppLookupAppfile(coreConfig.Appfile.File))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if app == nil {
		t.Fatal("app not found")
	}
	if app.Name != "foo" {
		t.Fatalf("bad: %#v", app)
	}
}

func TestCoreDev_compileMetadata(t *testing.T) {
	// Make a core that returns a fixed app
	coreConfig := TestCoreConfig(t)
	coreConfig.Appfile = TestAppfile(t, testPath("basic", "Appfile"))
	appMock := TestApp(t, TestAppTuple, coreConfig)
	core := testCore(t, coreConfig)

	// Configure the app to return a specific version in the metadata
	appMock.CompileResult = &app.CompileResult{Version: 12}

	// Compile!
	if err := core.Compile(); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Rebuild all the core so we have a fresh core
	core = testCore(t, coreConfig)

	// Run dev
	if err := core.Dev(); err != nil {
		t.Fatalf("err: %s", err)
	}

	// The context should have the right version
	if appMock.DevContext.CompileResult.Version != 12 {
		t.Fatalf("bad: %#v", appMock.DevContext.CompileResult)
	}
}

func testCore(t *testing.T, config *CoreConfig) *Core {
	core, err := NewCore(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return core
}
