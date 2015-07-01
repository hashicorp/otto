package otto

// Error is the interface implemented by many errors within Otto. You
// can use it to check what the type of an error is via the list of
// error codes below.
type Error interface {
	OriginalError() error
	Code() string
}

// codedError is the type used internally that implements the Error interface
type codedError struct {
	err  error
	code string
}

func (e *codedError) Error() string        { return e.err.Error() }
func (e *codedError) OriginalError() error { return e.err }
func (e *codedError) Code() string         { return e.code }

// errwrap.Wrapper impl.
func (e *codedError) WrappedErrors() []error { return []error{e.OriginalError()} }
