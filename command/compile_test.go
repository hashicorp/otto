package command

import (
	"testing"

	"github.com/mitchellh/cli"
)

func TestCompile(t *testing.T) {
	ui := new(cli.MockUi)
	c := &CompileCommand{
		Meta: Meta{
			Ui: ui,
		},
	}

	args := []string{}
	if code := c.Run(args); code != 0 {
		t.Fatalf("bad: %d\n\n%s", code, ui.ErrorWriter.String())
	}
}
