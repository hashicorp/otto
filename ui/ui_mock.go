package ui

// Mock is an implementation of Ui that stores its data in-memory
// primarily for testing purposes.
type Mock struct {
	HeaderBuf  []string
	MessageBuf []string
	RawBuf     []string

	InputCalled bool
	InputOpts   *InputOpts
	InputResult string
	InputError  error
}

func (u *Mock) Header(msg string) {
	u.HeaderBuf = append(u.HeaderBuf, msg)
}

func (u *Mock) Message(msg string) {
	u.MessageBuf = append(u.MessageBuf, msg)
}

func (u *Mock) Raw(msg string) {
	u.RawBuf = append(u.RawBuf, msg)
}

func (u *Mock) Input(opts *InputOpts) (string, error) {
	u.InputCalled = true
	u.InputOpts = opts
	return u.InputResult, u.InputError
}
