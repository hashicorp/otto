package command

import (
	"testing"

	"github.com/hashicorp/otto/ui"
)

func TestCliUi_impl(t *testing.T) {
	var _ ui.Ui = new(cliUi)
}
