package ui

import (
	"testing"
)

func TestNull_impl(t *testing.T) {
	var _ Ui = new(Null)
}
