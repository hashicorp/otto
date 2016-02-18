// This file contains the various tasks that are related to infras.

package otto

import (
	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/plan"
)

// TaskInfraCreds loads the infrastructure credentials. It will modify
// the context it points to in Extra but will also set the credentials
// in the result values.
type TaskInfraCreds struct {
	C *Core
}

func (t *TaskInfraCreds) Validate(args *plan.ExecArgs) (*plan.ExecResult, error) {
	return nil, nil
}

func (t *TaskInfraCreds) Execute(args *plan.ExecArgs) (*plan.ExecResult, error) {
	ctx := args.Extra["context"].(*context.Shared)

	// Get the infra implementation
	infra, infraCtx, err := t.C.infra()
	if err != nil {
		return nil, err
	}
	if err := t.C.infraCreds(infra, infraCtx); err != nil {
		return nil, err
	}
	defer maybeClose(infra)

	// Set it on our context
	ctx.InfraCreds = infraCtx.InfraCreds

	return nil, nil
}
