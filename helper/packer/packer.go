package packer

import (
	"os/exec"

	"github.com/hashicorp/otto/ui"
)

// Packer wraps `packer` execution into an easy-to-use API
type Packer struct {
	// Dir is the working directory where all Packer commands are executed
	Dir string

	// Ui, if given, will be used to stream output from the Packer commands.
	// If this is nil, then the output will be logged but won't be visible
	// to the user.
	Ui ui.Ui
}

// Execute executes a raw Packer command.
func (p *Packer) Execute(commandRaw ...string) error {
	// The command must always be machine-readable. We use this
	// exclusively to mirror the UI output.
	command := make([]string, len(commandRaw)+1)
	command[0] = "-machine-readable"
	copy(command[1:], commandRaw)

	// Build the command to execute
	cmd := exec.Command("packer", command...)
	cmd.Dir = p.Dir

	return nil
}
