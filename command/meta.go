package command

import (
	"github.com/mitchellh/cli"
)

// Meta are the meta-options that are available on all or most commands.
type Meta struct {
	Ui cli.Ui
}
