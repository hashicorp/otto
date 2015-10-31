package rpc

import (
	"reflect"
	"testing"

	"github.com/hashicorp/otto/ui"
)

func TestUi_impl(t *testing.T) {
	var _ ui.Ui = new(Ui)
}

func TestUi_input(t *testing.T) {
	client, server := testClientServer(t)
	defer client.Close()

	i := new(ui.Mock)
	i.InputResult = "foo"

	err := server.RegisterName("Ui", &UiServer{
		Ui: i,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	input := &Ui{Client: client, Name: "Ui"}

	opts := &ui.InputOpts{
		Id: "foo",
	}

	v, err := input.Input(opts)
	if !i.InputCalled {
		t.Fatal("input should be called")
	}
	if !reflect.DeepEqual(i.InputOpts, opts) {
		t.Fatalf("bad: %#v", i.InputOpts)
	}
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}

	if v != "foo" {
		t.Fatalf("bad: %#v", v)
	}
}
