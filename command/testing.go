package command

import (
	"testing"

	"github.com/hashicorp/otto/otto"
	"github.com/mitchellh/cli"
)

// TestMeta returns a Meta configured with in-memory structures.
func TestMeta(t *testing.T) Meta {
	core := otto.TestCoreConfig(t)
	ui := new(cli.MockUi)

	return Meta{
		CoreConfig: core,
		Ui:         ui,
	}
}
