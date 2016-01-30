// The plan package contains structures and helpers for Otto "plans,"
// the structure representing a goal and a set of tasks to achieve that
// goal.
package plan

import (
	"fmt"

	"github.com/hashicorp/otto/context"
)

// Plan is an executable object that represents a goal and the
// steps to take (tasks) to achieve that goal.
type Plan struct {
	Description string
	Tasks       []*Task
}

// Task is a single executable unit for a Plan. Tasks are meant to remain
// small in scope so that they can be composed and reasoned about within
// a plan.
type Task struct {
	Type string              // Type of the task
	Args map[string]*TaskArg // Args are the arguments to the Task

	Description         string // Short description of what this task will do
	DetailedDescription string // Long details about what this task will do (optional)
}

// TaskArg is an argument to a task. This is a struct rather than a plain
// "interface{}" to give us the option of adding more rich struct members
// later.
type TaskArg struct {
	// Value is the value of the argument
	Value interface{}
}

// TaskResult is a result type from executing a task.
//
// This is a struct rather than a raw "interface{}" so that we have the
// option of richer functions and struct members later.
type TaskResult struct {
	Value interface{}
}

// TaskExecutor is the interface that must be implemented to execute a
// task. The mapping of task "Type" to TaskExecutor is passed to Plan to
// execute.
type TaskExecutor interface {
	// Execute is called to run a task. It is given access to a populated
	// shared context and the map of arguments.
	Execute(ctx *context.Shared, args map[string]*TaskArg) (map[string]*TaskResult, error)
}

//-------------------------------------------------------------------
// GoStringer
//-------------------------------------------------------------------

func (p *Plan) GoString() string {
	return fmt.Sprintf("*%#v", *p)
}

func (t *Task) GoString() string {
	return fmt.Sprintf("*%#v", *t)
}
