package app

import (
	"testing"
)

func TestMock_impl(t *testing.T) {
	var _ App = new(Mock)
}
