package plan

type MockTaskExecutor struct {
	ValidateResult *ExecResult
	ValidateErr    error
}

func (e *MockTaskExecutor) Validate(*ExecArgs) (*ExecResult, error) {
	return e.ValidateResult, e.ValidateErr
}

func (e *MockTaskExecutor) Execute(*ExecArgs) (*ExecResult, error) {
	return nil, nil
}
