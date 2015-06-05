package ui

import (
	"testing"
)

func TestMock_impl(t *testing.T) {
	var _ Ui = new(Mock)
}
