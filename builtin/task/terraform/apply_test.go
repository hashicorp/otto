package terraform

import (
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/plan"
)

var hasTF = false

func init() {
	_, err := exec.LookPath("terraform")
	hasTF = err == nil
}

func TestApplyTask_impl(t *testing.T) {
	var _ plan.TaskExecutor = new(ApplyTask)
}

func TestApply(t *testing.T) {
	if !hasTF {
		t.Skip("Terraform not found")
	}

	var task ApplyTask

	// Build the args
	infraName := "foo"
	ctx := context.TestShared(t)
	pwd := filepath.Join("./test-fixtures", "basic")
	args := &plan.ExecArgs{
		Output: ioutil.Discard,
		Extra:  map[string]interface{}{"context": ctx},
		Args: map[string]*plan.TaskArg{
			"pwd":   &plan.TaskArg{Value: pwd},
			"infra": &plan.TaskArg{Value: infraName},
		},
	}

	// Initialize the infra
	lookup := &directory.InfraLookup{Name: infraName}
	infra := &directory.Infra{Name: infraName}
	if err := ctx.Directory.PutInfra(lookup, infra); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Validate
	if _, err := task.Validate(args); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Execute
	if _, err := task.Execute(args); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify the state was updated
	infra, err := ctx.Directory.GetInfra(lookup)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if len(infra.Opaque) <= 0 {
		t.Fatal("should have state")
	}
	if infra.State != directory.InfraStateReady {
		t.Fatal("state should be ready")
	}

	{
		expected := map[string]string{"foo": "bar"}
		if !reflect.DeepEqual(infra.Outputs, expected) {
			t.Fatalf("bad: %#v", infra.Outputs)
		}
	}
}
