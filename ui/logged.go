package ui

import (
	"log"
)

// Logged is an implementation of Ui that logs all messages as they
// pass through.
type Logged struct {
	Ui Ui
}

func (l *Logged) Header(msg string) {
	log.Printf("[INFO] ui header: %s", msg)
	l.Ui.Header(msg)
}

func (l *Logged) Message(msg string) {
	log.Printf("[INFO] ui message: %s", msg)
	l.Ui.Message(msg)
}

func (l *Logged) Raw(msg string) {
	log.Printf("[INFO] ui raw: %s", msg)
	l.Ui.Raw(msg)
}

func (l *Logged) Input(opts *InputOpts) (string, error) {
	// Not sure what to log here.
	return l.Ui.Input(opts)
}
