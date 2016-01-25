package infrastructure

import (
	"github.com/hashicorp/otto/plan"
)

// Mock is a mock implementation of the Infrastructure interface.
type Mock struct {
	CompileCalled  bool
	CompileContext *Context
	CompileResult  *CompileResult
	CompileErr     error

	PlanCalled  bool
	PlanContext *Context
	PlanResult  []*plan.Plan
	PlanErr     error
}

func (m *Mock) Creds(ctx *Context) (map[string]string, error) {
	return nil, nil
}

func (m *Mock) VerifyCreds(ctx *Context) error {
	return nil
}

func (m *Mock) Execute(ctx *Context) error {
	return nil
}

func (m *Mock) Compile(ctx *Context) (*CompileResult, error) {
	m.CompileCalled = true
	m.CompileContext = ctx
	return m.CompileResult, m.CompileErr
}

func (m *Mock) Plan(ctx *Context) ([]*plan.Plan, error) {
	m.PlanCalled = true
	m.PlanContext = ctx
	return m.PlanResult, m.PlanErr
}
