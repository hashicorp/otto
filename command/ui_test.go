package command

import (
	"bytes"
	"testing"

	"github.com/hashicorp/otto/ui"
)

func TestCliUi_impl(t *testing.T) {
	var _ ui.Ui = new(cliUi)
}

func TestCliUiInput(t *testing.T) {
	i := &cliUi{
		Reader: bytes.NewBufferString("foo\n"),
		Writer: bytes.NewBuffer(nil),
	}

	v, err := i.Input(&ui.InputOpts{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if v != "foo" {
		t.Fatalf("bad: %#v", v)
	}
}
