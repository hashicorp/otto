package otto

import (
	"testing"
)

func TestCodedError_impl(t *testing.T) {
	var _ Error = new(codedError)
}
