package otto

import (
	"testing"
)

func TestStyledUi_impl(t *testing.T) {
	var _ Ui = new(StyledUi)
}
