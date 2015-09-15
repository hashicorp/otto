package vagrant

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	execHelper "github.com/hashicorp/otto/helper/exec"
	"github.com/hashicorp/otto/ui"
)

//go:generate go-bindata -pkg=vagrant -nomemcopy -nometadata ./data/...

// Vagrant wraps `vagrant` execution into an easy-to-use API.
type Vagrant struct {
	// Dir is the working directory where all Vagrant commands will
	// be executed from.
	Dir string

	// DataDir is the directory where Vagrant commands should store data.
	DataDir string

	// Ui, if given, will be used to stream output from the Vagrant
	// commands. If this is nil, then the output will be logged but
	// won't be visible to the user.
	Ui ui.Ui
}

// A global mutex to prevent any Vagrant commands from running in parallel,
// which is not a supported mode of operation for Vagrant.
var vagrantMutex = &sync.Mutex{}

// The environment variable that Vagrant uses to configure its data dir.
const vagrantDataDirEnvVar = "VAGRANT_DOTFILE_PATH"

// Execute executes a raw Vagrant command.
func (v *Vagrant) Execute(command ...string) error {
	vagrantMutex.Lock()
	defer vagrantMutex.Unlock()

	// Build the command to execute
	cmd := exec.Command("vagrant", command...)
	cmd.Dir = v.Dir

	// Tell vagrant where to store data
	origDataDir := os.Getenv(vagrantDataDirEnvVar)
	defer os.Setenv(vagrantDataDirEnvVar, origDataDir)
	if err := os.Setenv(vagrantDataDirEnvVar, v.DataDir); err != nil {
		return err
	}

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
