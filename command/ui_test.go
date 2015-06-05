package command

import (
	"testing"

	"github.com/hashicorp/otto/otto"
)

func TestCliUi_impl(t *testing.T) {
	var _ otto.Ui = new(cliUi)
}
