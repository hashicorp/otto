package otto

import (
	"testing"

	"github.com/hashicorp/otto/plan"
)

func TestPlanExecute_validate(t *testing.T) {
	mock := &plan.MockTask{}
	opts := &PlanOpts{
		Validate: true,
		extraTasks: map[string]plan.TaskExecutor{
			"test": mock,
		},
	}

	p := &Plan{Plans: plan.TestPlan(t, testPath("plan-basic", "infra.hcl"))}
	err := p.Execute(TestCore(t, nil), opts)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !mock.ValidateCalled {
		t.Fatal("should call validate")
	}
	if mock.ExecuteCalled {
		t.Fatal("should not call execute")
	}
}

func TestPlanExecute_execute(t *testing.T) {
	mock := &plan.MockTask{}
	opts := &PlanOpts{
		extraTasks: map[string]plan.TaskExecutor{
			"test": mock,
		},
	}

	p := &Plan{Plans: plan.TestPlan(t, testPath("plan-basic", "infra.hcl"))}
	err := p.Execute(TestCore(t, nil), opts)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if mock.ValidateCalled {
		t.Fatal("should call validate")
	}
	if !mock.ExecuteCalled {
		t.Fatal("should not call execute")
	}
}
