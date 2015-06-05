package command

import (
	"fmt"

	"github.com/hashicorp/otto/ui"
	"github.com/mitchellh/cli"
)

// NewUi returns a new otto Ui implementation for use around
// the given CLI Ui implementation.
func NewUi(raw cli.Ui) ui.Ui {
	return &ui.Styled{
		Ui: &cliUi{
			CliUi: raw,
		},
	}
}

// cliUi is a wrapper around a cli.Ui that implements the otto.Ui
// interface. It is unexported since the NewUi method should be used
// instead.
type cliUi struct {
	CliUi cli.Ui
}

func (u *cliUi) Header(msg string) {
	u.CliUi.Output(msg)
}

func (u *cliUi) Message(msg string) {
	u.CliUi.Output(msg)
}

func (u *cliUi) Raw(msg string) {
	fmt.Print(msg)
}
