package otto

// StyledUi is a wrapper around an existing UI that automatically
// adds formatting around the UI text.
type StyledUi struct {
	Ui
}

func (u *StyledUi) Header(msg string) {
	u.Ui.Header("==> " + msg)
}

func (u *StyledUi) Message(msg string) {
	u.Ui.Message("    " + msg)
}
