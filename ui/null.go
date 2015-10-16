package ui

import (
	"errors"
)

// Null is an implementation of Ui that does nothing.
type Null struct{}

func (*Null) Header(string)  {}
func (*Null) Message(string) {}
func (*Null) Raw(string)     {}
func (*Null) Input(*InputOpts) (string, error) {
	return "", errors.New("null ui")
}
