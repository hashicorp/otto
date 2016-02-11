// The plan package contains structures and helpers for Otto "plans,"
// the structure representing a goal and a set of tasks to achieve that
// goal.
package plan

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hil"
	"github.com/hashicorp/hil/ast"
)

// Plan is an executable object that represents a goal and the
// steps to take (tasks) to achieve that goal.
type Plan struct {
	Description string  // Description of what the plan does
	Tasks       []*Task // Serial order of tasks to execute

	// Inputs is a map of variables to set as inputs to the plan.
	// These are available as "input.NAME" within the tasks.
	Inputs map[string]interface{}
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

//-------------------------------------------------------------------
// Methods on structs
//-------------------------------------------------------------------

// Interpolate does the interpolation on the value and returns a new
// TaskArg copy that has the interpolated value. If the value is not
// interpolated, the arg returned is not a copy.
func (t *TaskArg) Interpolate(vs map[string]*TaskResult) (*TaskArg, error) {
	// Parse
	root, err := t.hil()
	if err != nil {
		return nil, err
	}
	if root == nil {
		return t, nil
	}

	// Build the variables
	varMap := make(map[string]ast.Variable)
	for k, raw := range vs {
		v, err := ast.NewVariable(raw.Value)
		if err != nil {
			return nil, fmt.Errorf("var %s: %s", k, err)
		}

		varMap[k] = v
	}

	// Eval
	v, _, err := hil.Eval(root, &hil.EvalConfig{
		GlobalScope: &ast.BasicScope{
			VarMap: varMap,
		},
	})
	if err != nil {
		return nil, err
	}

	// Return new arg
	return &TaskArg{
		Value: v,
	}, nil
}

// Refs returns the references to other variables found within the
// argument.
func (t *TaskArg) Refs() []string {
	// Parse the value as HIL
	root, err := t.hil()
	if err != nil {
		// Panic, this should've been validated earlier
		panic(err)
	}
	if root == nil {
		return nil
	}

	// Walk the AST and find all variable access
	varMap := make(map[string]struct{})
	root.Accept(func(n ast.Node) ast.Node {
		if vn, ok := n.(*ast.VariableAccess); ok {
			varMap[vn.Name] = struct{}{}
		}

		return n
	})

	// Create the result slice from the map so that all values are unique
	result := make([]string, 0, len(varMap))
	for k, _ := range varMap {
		result = append(result, k)
	}

	sort.Strings(result)
	return result
}

func (t *TaskArg) hil() (ast.Node, error) {
	s, ok := t.Value.(string)
	if !ok {
		// If it isn't a string it can't reference any other vars
		return nil, nil
	}

	// Parse the value as HIL
	return hil.Parse(s)
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
