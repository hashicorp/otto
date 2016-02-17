package terraform

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/directory"
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
	infraName := args.Args["infra"].Value.(string)
	ctx := args.Extra["context"].(*context.Shared)

	// Temporary directory for scratch
	tempDir, err := ioutil.TempDir("", "terraform")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tempDir)

	// State paths will always be here
	state := filepath.Join(tempDir, "state.in")
	stateOut := filepath.Join(tempDir, "state.out")

	// Get the infrastructure if we can
	lookup := &directory.InfraLookup{Name: infraName}
	infra, err := ctx.Directory.GetInfra(lookup)
	if err != nil {
		return nil, err
	}
	if infra == nil {
		// If we have no infrastructure then we also consider it an error
		return nil, fmt.Errorf("infra not found in directory: %s", infraName)
	}

	// If we have state, write it out
	if len(infra.Opaque) > 0 {
		if err := ioutil.WriteFile(state, infra.Opaque, 0644); err != nil {
			return nil, err
		}
	}

	// Variables for Terraform will go in here as we build them
	vars := make(map[string]string)

	// If we have credentials, setup those args
	if len(ctx.InfraCreds) > 0 {
		for k, v := range ctx.InfraCreds {
			vars[k] = v
		}
	}

	// Build the execArgs for Terraform here
	execArgs := make([]string, 0, 50)
	execArgs = append(execArgs, "apply")
	execArgs = append(execArgs, "-state", state)
	execArgs = append(execArgs, "-state-out", stateOut)
	for k, v := range vars {
		execArgs = append(execArgs, "-var", fmt.Sprintf("%s=%s", k, v))
	}

	// Build the command. Terraform is idempotent so we can just
	// run apply as-is multiple times if this task is repeated.
	var output bytes.Buffer
	stdout := io.MultiWriter(args.Output, &output)
	cmd := exec.Command("terraform", execArgs...)
	cmd.Dir = pwd
	cmd.Stdin = os.Stdin
	cmd.Stdout = stdout
	cmd.Stderr = stdout

	// No return value, no error!
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("Error executing Terraform: %s\n\n"+
			"The raw output from Terraform is below:\n\n%s", err, output.String())
		output.Reset()
	}

	// Update the state of the infra
	infra.State = directory.InfraStateReady
	if err != nil {
		infra.State = directory.InfraStatePartial
	}

	// Always store the output state
	if _, err := os.Stat(stateOut); err == nil {
		infra.Opaque, err = ioutil.ReadFile(stateOut)
		if err != nil {
			return nil, fmt.Errorf(errStateSaveFailed)
		}

		infra.Outputs, err = Outputs(stateOut)
		if err != nil {
			return nil, fmt.Errorf(errStateSaveFailed)
		}
	}

	// Update the infra
	if e := ctx.Directory.PutInfra(lookup, infra); e != nil {
		err = e
	}

	// Return
	return nil, err
}

const errStateSaveFailed = `Failed to save Terraform state: %s

This means that Otto was unable to store the state of your infrastructure.
At this time, Otto doesn't support gracefully recovering from this
scenario. The state should be in the path below. Please ask the community
for assistance, and do not lose the file below.`
