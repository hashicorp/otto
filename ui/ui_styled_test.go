package ui

import (
	"testing"
)

func TestStyled_impl(t *testing.T) {
	var _ Ui = new(Styled)
}
