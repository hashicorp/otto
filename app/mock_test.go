package app

import (
	"io"
	"testing"
)

func TestMock_impl(t *testing.T) {
	var _ App = new(Mock)
	var _ io.Closer = new(Mock)
}
