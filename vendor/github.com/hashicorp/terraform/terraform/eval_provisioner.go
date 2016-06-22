package terraform

import (
	"fmt"
)

// EvalInitProvisioner is an EvalNode implementation that initializes a provisioner
// and returns nothing. The provisioner can be retrieved again with the
// EvalGetProvisioner node.
type EvalInitProvisioner struct {
	Name string
}

func (n *EvalInitProvisioner) Eval(ctx EvalContext) (interface{}, error) {
	return ctx.InitProvisioner(n.Name)
}

// EvalCloseProvisioner is an EvalNode implementation that closes provisioner
// connections that aren't needed anymore.
type EvalCloseProvisioner struct {
	Name string
}

func (n *EvalCloseProvisioner) Eval(ctx EvalContext) (interface{}, error) {
	ctx.CloseProvisioner(n.Name)
	return nil, nil
}

// EvalGetProvisioner is an EvalNode implementation that retrieves an already
// initialized provisioner instance for the given name.
type EvalGetProvisioner struct {
	Name   string
	Output *ResourceProvisioner
}

func (n *EvalGetProvisioner) Eval(ctx EvalContext) (interface{}, error) {
	result := ctx.Provisioner(n.Name)
	if result == nil {
		return nil, fmt.Errorf("provisioner %s not initialized", n.Name)
	}

	if n.Output != nil {
		*n.Output = result
	}

	return result, nil
}
