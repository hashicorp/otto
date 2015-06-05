package otto

// UiMock is an implementation of Ui that stores its data in-memory
// primarily for testing purposes.
type UiMock struct {
	HeaderBuf  []string
	MessageBuf []string
	RawBuf     []string
}

func (u *UiMock) Header(msg string) {
	u.HeaderBuf = append(u.HeaderBuf, msg)
}

func (u *UiMock) Message(msg string) {
	u.MessageBuf = append(u.MessageBuf, msg)
}

func (u *UiMock) Raw(msg string) {
	u.RawBuf = append(u.RawBuf, msg)
}
