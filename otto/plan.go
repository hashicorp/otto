package otto

import (
	"github.com/hashicorp/otto/plan"
)

// Plan is a deployment plan for Otto.
type Plan struct {
	// Infra is the list of low-level plans the infra has
	Infra []*plan.Plan
}

// PlanOpts are options for plan execution.
type PlanOpts struct {
	// Validate, if true, will only validate the plan and not execute it.
	Validate bool

	// Used for testing only: specifies extra tasks to inject
	extraTasks map[string]plan.TaskExecutor
}

// TODO: test
func (p *Plan) Empty() bool {
	return p != nil && len(p.Infra) == 0
}

// Execute will execute the plan. Depending on the PlanOpts, different
// parts of the plan may be executed.
func (p *Plan) Execute(c *Core, opts *PlanOpts) error {
	// Get the task map from our core
	taskMap, err := c.planTaskMap()
	if err != nil {
		return err
	}
	for k, v := range opts.extraTasks {
		taskMap[k] = v
	}

	// Make sure we close all tasks if they support it
	for _, t := range taskMap {
		defer maybeClose(t)
	}

	// Instantiate the plan executor
	e := &plan.Executor{TaskMap: taskMap}

	// Get the function we need to call
	var f func(*plan.Plan) error = e.Validate
	if !opts.Validate {
		f = e.Execute
	}

	// Go through the plans in the proper order: infra, foundation, app
	for _, p := range p.Infra {
		if err := f(p); err != nil {
			return err
		}
	}

	return nil
}

//--------------------------------------------------------------------
// Core Methods
//--------------------------------------------------------------------

func (c *Core) planTaskMap() (map[string]plan.TaskExecutor, error) {
	result := make(map[string]plan.TaskExecutor)

	// Static built-ins
	result["delete"] = &plan.DeleteTask{}
	result["store"] = &plan.StoreTask{}

	// From plugins
	// TODO

	return result, nil
}
