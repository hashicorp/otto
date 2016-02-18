package otto

import (
	"reflect"
	"testing"

	"github.com/hashicorp/otto/ui"
)

func TestCoreInfraCreds(t *testing.T) {
	// We've never asked for infra creds before.

	// Prepare the config
	uiMock := new(ui.Mock)
	config := TestCoreConfig(t)
	config.Ui = &ui.Logged{Ui: uiMock}
	infra := TestInfra(t, "test", config)
	core, err := NewCore(config)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Setup what we'll return
	infra.CredsResult = map[string]string{"foo": "bar"}

	// Setup our password value
	uiMock.InputResult = "foo"

	// Get the infra
	infraImpl, infraCtx, err := core.infra()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Get creds
	if err := core.infraCreds(infraImpl, infraCtx); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify they're correct
	if !reflect.DeepEqual(infraCtx.InfraCreds, infra.CredsResult) {
		t.Fatalf("bad: %#v", infraCtx.InfraCreds)
	}

	// Reset some things
	uiMock.InputCalled = false
	infra.CredsCalled = false
	infraCtx.InfraCreds = nil

	// Getting it again, test that we don't call Creds this time
	// since it should be cached.
	if err := core.infraCreds(infraImpl, infraCtx); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify
	if infra.CredsCalled {
		t.Fatal("creds shouldn't be called again")
	}
	if !reflect.DeepEqual(infraCtx.InfraCreds, infra.CredsResult) {
		t.Fatalf("bad: %#v", infraCtx.InfraCreds)
	}
}
