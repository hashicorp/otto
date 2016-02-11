package terraform

import (
	"os"
	"os/exec"

	"github.com/hashicorp/otto/plan"
)

// ApplyTask is a plan.TaskExecutor that runs Terraform.
//
// Args:
//
//   - pwd (string) - working directory to execute from
//   - state (string) - path to the state file
//   - state_out (string) - path to the output state file
//
// Execution:
//
// Even if this task errors, if the state_out path exists, then it should
// be considered a partial state and stored.
//
type ApplyTask struct{}

func (t *ApplyTask) Validate(args *plan.ExecArgs) (*plan.ExecResult, error) {
	return nil, nil
}

func (t *ApplyTask) Execute(args *plan.ExecArgs) (*plan.ExecResult, error) {
	// Get all our args
	pwd := args.Args["pwd"].Value.(string)
	state := args.Args["state"].Value.(string)
	stateOut := args.Args["state_out"].Value.(string)

	// Build the command. Terraform is idempotent so we can just
	// run apply as-is multiple times if this task is repeated.
	cmd := exec.Command(
		"terraform",
		"apply",
		"-state", state,
		"-state-out", stateOut,
	)
	cmd.Dir = pwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = args.Output
	cmd.Stderr = args.Output

	// No return value, no error!
	return nil, cmd.Run()
}
