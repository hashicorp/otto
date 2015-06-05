package command

import (
	"testing"

	"github.com/hashicorp/otto/otto"
	"github.com/mitchellh/cli"
)

func TestCompile(t *testing.T) {
	core := otto.TestCoreConfig(t)
	infra := otto.TestInfra(t, "test", core)
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: core,
			Ui:         ui,
		},
	}

	args := []string{"-appfile", fixtureDir("compile-basic")}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}

	if !infra.CompileCalled {
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

	args := []string{"-appfile", fixtureDir("compile-file/Appfile.other")}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
