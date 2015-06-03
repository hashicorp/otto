package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestCompile(t *testing.T) {
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: testCoreConfig(t),
			Ui:         ui,
		},
	}

	args := []string{fixtureDir("compile-basic")}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}

func TestCompile_pathFile(t *testing.T) {
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			CoreConfig: testCoreConfig(t),
			Ui:         ui,
		},
	}

	args := []string{fixtureDir("compile-file/Appfile.other")}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
