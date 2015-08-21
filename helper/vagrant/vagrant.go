package vagrant

import (
	"fmt"
	"os/exec"

	execHelper "github.com/hashicorp/otto/helper/exec"
	"github.com/hashicorp/otto/ui"
)

// Vagrant wraps `vagrant` execution into an easy-to-use API.
type Vagrant struct {
	// Dir is the working directory where all Vagrant commands will
	// be executed from.
	Dir string

	// Ui, if given, will be used to stream output from the Vagrant
	// commands. If this is nil, then the output will be logged but
	// won't be visible to the user.
	Ui ui.Ui
}

// Execute executes a raw Vagrant command.
func (v *Vagrant) Execute(command ...string) error {
	// Build the command to execute
	cmd := exec.Command("vagrant", command...)
	cmd.Dir = v.Dir

	// Run it with the execHelper
	if err := execHelper.Run(v.Ui, cmd); err != nil {
		return fmt.Errorf(
			"Error executing Vagrant: %s\n\n"+
				"The error messages from Vagrant are usually very informative.\n"+
				"Please read it carefully and fix any issues it mentions. If\n"+
				"the message isn't clear, please report this to the Otto project.",
			err)
	}

	return nil
}
