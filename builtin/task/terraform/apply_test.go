package terraform

import (
	"testing"

	"github.com/hashicorp/otto/plan"
)

func TestApplyTask_impl(t *testing.T) {
	var _ plan.TaskExecutor = new(ApplyTask)
}
