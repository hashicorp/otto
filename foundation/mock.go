package foundation

// Mock is a mock implementation of the Foundation interface.
type Mock struct {
	CompileCalled  bool
	CompileContext *Context
	CompileResult  *CompileResult
	CompileErr     error

	InfraCalled  bool
	InfraContext *Context
	InfraErr     error
}

func (m *Mock) Compile(ctx *Context) (*CompileResult, error) {
	m.CompileCalled = true
	m.CompileContext = ctx
	return m.CompileResult, m.CompileErr
}

func (m *Mock) Infra(ctx *Context) error {
	m.InfraCalled = true
	m.InfraContext = ctx
	return m.InfraErr
}
