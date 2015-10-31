package otto

import (
	"testing"

	"github.com/hashicorp/otto/app"
)

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
