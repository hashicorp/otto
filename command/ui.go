package command

import (
	"github.com/hashicorp/otto/otto"
	"github.com/mitchellh/cli"
)

// NewUi returns a new otto Ui implementation for use around
// the given CLI Ui implementation.
func NewUi(ui cli.Ui) otto.Ui {
	return &otto.StyledUi{
		Ui: &cliUi{
			CliUi: ui,
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
	u.CliUi.Output(msg)
}
