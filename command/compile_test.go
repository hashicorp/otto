package command

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/otto"
	"github.com/mitchellh/cli"
)

func TestCompile(t *testing.T) {
	core := otto.TestCoreConfig(t)
	infra := otto.TestInfra(t, "test", core)
	appImpl := otto.TestApp(t, app.Tuple{"test", "test", "test"}, core)
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: core,
			Ui:         ui,
		},
	}

	dir := fixtureDir("compile-basic")
	defer os.Remove(filepath.Join(dir, ".ottoid"))

	args := []string{"-appfile", dir}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	if !infra.CompileCalled {
		t.Fatal("Compile should be called")
	}
	if !appImpl.CompileCalled {
		t.Fatal("Compile should be called")
	}
}

func TestCompile_pathFile(t *testing.T) {
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: otto.TestCoreConfig(t),
			Ui:         ui,
		},
	}

	dir := fixtureDir("compile-file")
	defer os.Remove(filepath.Join(dir, ".ottoid"))

	args := []string{"-appfile", filepath.Join(dir, "Appfile.other")}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
