package plan

import (
	"testing"
)

func TestMockTask_impl(t *testing.T) {
	var _ TaskExecutor = new(MockTask)
}
