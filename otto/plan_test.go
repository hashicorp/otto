package otto

import (
	"testing"

	"github.com/hashicorp/otto/plan"
	"github.com/hashicorp/otto/ui"
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

func TestPlanExecute_executeUi(t *testing.T) {
	// Setup mock to output
	mock := &plan.MockTask{}
	mock.ExecuteFn = func(x *plan.ExecArgs) (*plan.ExecResult, error) {
		x.Println("HELLO!")
		return nil, nil
	}

	// Configure the task
	opts := &PlanOpts{
		extraTasks: map[string]plan.TaskExecutor{
			"test": mock,
		},
	}

	// Configure the core to have a UI
	ui := &ui.Mock{}
	core := TestCore(t, nil)
	core.ui = ui

	p := &Plan{Plans: plan.TestPlan(t, testPath("plan-basic", "infra.hcl"))}
	err := p.Execute(core, opts)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify we saw it
	if len(ui.RawBuf) == 0 {
		t.Fatal("should have raw messages")
	}
	if ui.RawBuf[0] != "HELLO!\n" {
		t.Fatalf("bad: %#v", ui.RawBuf)
	}
}
