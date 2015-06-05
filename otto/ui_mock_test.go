package otto

import (
	"testing"
)

func TestMockUi_impl(t *testing.T) {
	var _ Ui = new(UiMock)
}
