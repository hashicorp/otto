package ottoInfra

import (
	"io/ioutil"
	"testing"

	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/plan"
)

func TestDeployVersionTask_impl(t *testing.T) {
	var _ plan.TaskExecutor = new(DeployVersionTask)
}

func TestDeployVersion(t *testing.T) {
	var task DeployVersionTask

	// Build the args
	infraName := "foo"
	ctx := context.TestShared(t)
	args := &plan.ExecArgs{
		Output: ioutil.Discard,
		Extra:  map[string]interface{}{"context": ctx},
		Args: map[string]*plan.TaskArg{
			"infra":          &plan.TaskArg{Value: infraName},
			"deploy_version": &plan.TaskArg{Value: "4.2.3"},
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
	if infra.DeployVersion != "4.2.3" {
		t.Fatalf("bad: %#v", infra.DeployVersion)
	}
}
