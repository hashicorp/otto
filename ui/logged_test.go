package ui

import (
	"testing"
)

func TestLogged_impl(t *testing.T) {
	var _ Ui = new(Logged)
}
