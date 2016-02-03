package plan

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/otto/context"
)

// TaskExecutor is the interface that must be implemented to execute a
// task. The mapping of task "Type" to TaskExecutor is passed to Plan to
// execute.
type TaskExecutor interface {
	// Validate is called to validate the arguments and hint at return values.
	Validate(*ExecArgs) (*ExecResult, error)

	// Execute is called to perform the actual task.
	Execute(*ExecArgs) (*ExecResult, error)
}

// ExecArgs are the arguments given to a TaskExecutor.
type ExecArgs struct {
	// Ctx is the Otto context for this execution
	Ctx *context.Shared

	// Args is the map of arguments and their value. For validation,
	// the TaskArg value will be uninterpolated and thus shouldn't be
	// used. Keys can be used for validation.
	Args map[string]*TaskArg
}

// ExecResult is the result returned from a TaskExecutor
type ExecResult struct {
	// Values are the resulting named values from the execution.
	Values map[string]*TaskResult

	// Store can be used to put values into storage. This shouldn't be used
	// publicly. It is exposed in case it MUST be used but this is meant
	// to only be used by the "Store" task type. If *TaskResult is nil, then
	// it will be deleted from the store.
	Store map[string]*TaskResult
}

// Executor is the struct used to execute a plan.
type Executor struct {
	// Callback, if non-nil, will be called for various events during
	// execution. You can use this to get information and control the
	// execution.
	Callback func(ExecuteEvent)

	// TaskMap is the map of Task types to executors for that task
	TaskMap map[string]TaskExecutor
}

// Validate will validate the semantics of the plan. This checks that
// all variable access will resolve, all task types are valid, etc.
func (e *Executor) Validate(p *Plan, ctx *context.Shared) error {
	var err error

	// First verify all the task types are valid
	for _, t := range p.Tasks {
		if _, ok := e.TaskMap[t.Type]; !ok {
			err = multierror.Append(err, fmt.Errorf("Unknown task type: %s", t.Type))
		}
	}

	// If we have errors at this point just return since the rest of the
	// checks will be difficult.
	if err != nil {
		return err
	}

	// Now go through all the tasks and validate the arg keys and the
	// variable access. The varMap below wil keep track of the variables
	// we'll have.
	varMap := make(map[string]struct{})
	resultMap := make(map[string]struct{})
	for i, t := range p.Tasks {
		// Create the full map of available vars
		fullMap := make(map[string]struct{})
		for k, v := range resultMap {
			fullMap[fmt.Sprintf("result.%s", k)] = v
		}
		for k, v := range varMap {
			fullMap[k] = v
		}

		// Validate the vars in the args
		for _, a := range t.Args {
			for _, ref := range a.Refs() {
				if _, ok := fullMap[ref]; !ok {
					err = multierror.Append(err, fmt.Errorf(
						"Task %d (%s): unknown reference: %s", i+1, t.Type, ref))
				}
			}
		}

		// Call Validate to validate the args
		te := e.TaskMap[t.Type]
		args := &ExecArgs{Ctx: ctx, Args: t.Args}
		result, verr := te.Validate(args)
		if verr != nil {
			err = multierror.Append(err, multierror.Prefix(
				verr, fmt.Sprintf("Task %d (%s): ", i+1, t.Type)))
			break
		}

		// Keep track of the result types
		resultMap = make(map[string]struct{})
		for k, _ := range result.Values {
			resultMap[k] = struct{}{}
		}

		// Keep track of storage
		for k, v := range result.Store {
			if v == nil {
				delete(varMap, k)
			} else {
				varMap[k] = struct{}{}
			}
		}
	}

	return err
}

// Execute is called to execute a plan.
//
// The configured Callback mechanism can be used to get regular progress
// events and control the execution. This function will block.
func (e *Executor) Execute(p *Plan, ctx *context.Shared) error {
	// Execute the tasks in serial
	for _, t := range p.Tasks {
		te := e.TaskMap[t.Type]
		args := &ExecArgs{Ctx: ctx, Args: t.Args}
		if _, err := te.Execute(args); err != nil {
			return err
		}
	}

	return nil
}

// ExecuteEvent is an event that a callback can receive during execution.
// You must type switch on the various implementations below.
type ExecuteEvent interface{}
