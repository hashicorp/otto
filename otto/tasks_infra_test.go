package otto

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/plan"
	"github.com/hashicorp/otto/ui"
)

func TestTaskInfraCreds_impl(t *testing.T) {
	var _ plan.TaskExecutor = new(TaskInfraCreds)
}

func TestTaskInfraCreds(t *testing.T) {
	expected := map[string]string{"foo": "bar"}

	infra := new(infrastructure.Mock)
	uiMock := new(ui.Mock)
	core := TestCore(t, &TestCoreOpts{Infra: infra})
	core.ui = uiMock
	task := &TaskInfraCreds{C: core}

	// Return value for creds
	infra.CredsResult = expected

	// Return value so we can answer questions
	uiMock.InputResult = "foo"

	// Build the args
	ctx := context.TestShared(t)
	args := &plan.ExecArgs{
		Output: ioutil.Discard,
		Extra:  map[string]interface{}{"context": ctx},
	}

	// Validate
	if _, err := task.Validate(args); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Execute
	if _, err := task.Execute(args); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify that the credentials are set on the context
	if !reflect.DeepEqual(ctx.InfraCreds, expected) {
		t.Fatalf("bad: %#v", ctx.InfraCreds)
	}
}
