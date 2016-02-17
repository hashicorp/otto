package ottoInfra

import (
	"fmt"

	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/plan"
)

// DeployVersionTask is a plan.TaskExecutor that updates the deploy
// version of an infrastructure.
//
// Args:
//
//   - infra (string) - The name of the infra
//   - deploy_version (string) - The version to update to
//
type DeployVersionTask struct{}

func (t *DeployVersionTask) Validate(args *plan.ExecArgs) (*plan.ExecResult, error) {
	return nil, nil
}

func (t *DeployVersionTask) Execute(args *plan.ExecArgs) (*plan.ExecResult, error) {
	// Get all our args
	infraName := args.Args["infra"].Value.(string)
	deployVersion := args.Args["deploy_version"].Value.(string)
	ctx := args.Extra["context"].(*context.Shared)

	// Get the infrastructure if we can
	lookup := &directory.InfraLookup{Name: infraName}
	infra, err := ctx.Directory.GetInfra(lookup)
	if err != nil {
		return nil, err
	}
	if infra == nil {
		// If we have no infrastructure then we also consider it an error
		return nil, fmt.Errorf("infra not found in directory: %s", infraName)
	}

	// Update the version
	infra.DeployVersion = deployVersion

	// Update the infra
	if e := ctx.Directory.PutInfra(lookup, infra); e != nil {
		err = e
	}

	// Return
	return nil, err
}
