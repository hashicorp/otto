package plan

import (
	"fmt"
)

// DeleteTask is a built-in type of task that is used to store value
// in the "memory" of plan execution. These can then be referenced
// directly with $foo where "foo" is the name of the variable.
type DeleteTask struct{}

func (t *DeleteTask) Validate(args *ExecArgs) (*ExecResult, error) {
	if len(args.Args) != 1 {
		return nil, fmt.Errorf("exactly one arg 'key' should be given")
	}

	arg, ok := args.Args["key"]
	if !ok {
		return nil, fmt.Errorf("exactly one arg 'key' should be given")
	}

	if len(arg.Refs()) > 0 {
		return nil, fmt.Errorf("key can't contain interpolations")
	}

	value, ok := arg.Value.(string)
	if !ok {
		return nil, fmt.Errorf("key must be a string")
	}

	resultMap := make(map[string]*TaskResult)
	resultMap[value] = nil
	return &ExecResult{Store: resultMap}, nil
}

func (t *DeleteTask) Execute(args *ExecArgs) (*ExecResult, error) {
	// TODO
	return nil, nil
}
