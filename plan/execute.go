package plan

import (
	"github.com/hashicorp/otto/context"
)

// Executor is the struct used to execute a plan.
type Executor struct {
	// Callback, if non-nil, will be called for various events during
	// execution. You can use this to get information and control the
	// execution.
	Callback func(ExecuteEvent)

	// TaskMap is the map of Task types to executors for that task
	TaskMap map[string]TaskExecutor
}

// Execute is called to execute a plan.
//
// The configured Callback mechanism can be used to get regular progress
// events and control the execution. This function will block.
func (e *Executor) Execute(p *Plan, ctx *context.Shared) error {
	// Execute the tasks in serial
	for _, t := range p.Tasks {
		te := e.TaskMap[t.Type]
		if _, err := te.Execute(ctx, t.Args); err != nil {
			return err
		}
	}

	return nil
}

// ExecuteEvent is an event that a callback can receive during execution.
// You must type switch on the various implementations below.
type ExecuteEvent interface{}
