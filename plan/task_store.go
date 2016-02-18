package plan

// StoreTask is a built-in type of task that is used to store value
// in the "memory" of plan execution. These can then be referenced
// directly with $foo where "foo" is the name of the variable.
type StoreTask struct{}

func (t *StoreTask) Validate(args *ExecArgs) (*ExecResult, error) {
	resultMap := make(map[string]*TaskResult, len(args.Args))
	for k, _ := range args.Args {
		resultMap[k] = &TaskResult{Value: ""}
	}

	return &ExecResult{Store: resultMap}, nil
}

func (t *StoreTask) Execute(args *ExecArgs) (*ExecResult, error) {
	store := make(map[string]*TaskResult)
	for k, a := range args.Args {
		store[k] = &TaskResult{
			Value: a.Value.(string),
		}
	}

	return &ExecResult{Store: store}, nil
}
