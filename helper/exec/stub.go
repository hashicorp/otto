package exec

import (
	"os/exec"
)

// Runner is the function that is called to run a command. This can be
// overridden for tests. This is a low-level replacement but we use it in
// various helper libraries for unit testing.
var Runner func(*exec.Cmd) error = RealRunner

// TestChrunner is a helper function that can be used to temporarily change
// the runner. This returns a function call that should be defered to restore
// the runner. Example:
//
//     defer TestChrunner(newRunner)()
//
func TestChrunner(r func(*exec.Cmd) error) func() {
	oldRunner := Runner
	Runner = r
	return func() {
		Runner = oldRunner
	}
}

// RealRunner is the default value of Runner and actually executes a command.
func RealRunner(cmd *exec.Cmd) error {
	return cmd.Run()
}

// MockRunner is a Runner implementation that records the calls. The
// Runner can be set to the MockRunner's Run function.
type MockRunner struct {
	// This will record the commands that are executed
	Commands []*exec.Cmd

	// This will be the return values of the errors for the commands
	// that are executed, in the order they're called. If this is empty
	// or shorter than the command, a nil error is returned.
	CommandErrs []error
}

func (r *MockRunner) Run(cmd *exec.Cmd) error {
	r.Commands = append(r.Commands, cmd)
	if len(r.CommandErrs) < len(r.Commands) {
		return nil
	}

	return r.CommandErrs[len(r.Commands)-1]
}
