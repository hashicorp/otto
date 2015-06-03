package infrastructure

import (
	"testing"
)

func TestMock_impl(t *testing.T) {
	var _ Infrastructure = new(Mock)
}
