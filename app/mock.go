package app

// Mock is a mock implementation of the App interface.
type Mock struct {
	CompileCalled  bool
	CompileContext *Context
	CompileResult  *CompileResult
	CompileErr     error
	CompileFunc    func(ctx *Context) (*CompileResult, error)

	BuildCalled  bool
	BuildContext *Context
	BuildErr     error

	DeployCalled  bool
	DeployContext *Context
	DeployErr     error

	DevCalled  bool
	DevContext *Context
	DevErr     error

	DevDepCalled     bool
	DevDepContextDst *Context
	DevDepContextSrc *Context
	DevDepResult     *DevDep
	DevDepErr        error
}

func (m *Mock) Compile(ctx *Context) (*CompileResult, error) {
	m.CompileCalled = true
	m.CompileContext = ctx
	if m.CompileFunc != nil {
		return m.CompileFunc(ctx)
	}
	return m.CompileResult, m.CompileErr
}

func (m *Mock) Build(ctx *Context) error {
	m.BuildCalled = true
	m.BuildContext = ctx
	return m.BuildErr
}

func (m *Mock) Deploy(ctx *Context) error {
	m.DeployCalled = true
	m.DeployContext = ctx
	return m.DeployErr
}

func (m *Mock) Dev(ctx *Context) error {
	m.DevCalled = true
	m.DevContext = ctx
	return m.DevErr
}

func (m *Mock) DevDep(dst, src *Context) (*DevDep, error) {
	m.DevDepCalled = true
	m.DevDepContextDst = dst
	m.DevDepContextSrc = src
	return m.DevDepResult, m.DevDepErr
}
