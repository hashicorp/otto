package plan

// MockTask implements TaskExecutor and can be used for testing.
type MockTask struct {
	ValidateCalled bool
	ValidateArgs   *ExecArgs
	ValidateResult *ExecResult
	ValidateErr    error

	ExecuteCalled bool
	ExecuteArgs   *ExecArgs
	ExecuteResult *ExecResult
	ExecuteErr    error
}

func (t *MockTask) Validate(args *ExecArgs) (*ExecResult, error) {
	t.ValidateCalled = true
	t.ValidateArgs = args
	return t.ValidateResult, t.ValidateErr
}

func (t *MockTask) Execute(args *ExecArgs) (*ExecResult, error) {
	t.ExecuteCalled = true
	t.ExecuteArgs = args
	return t.ExecuteResult, t.ExecuteErr
}
