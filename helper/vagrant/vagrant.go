package vagrant

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/otto/context"
	execHelper "github.com/hashicorp/otto/helper/exec"
	"github.com/hashicorp/otto/helper/hashitools"
	"github.com/hashicorp/otto/ui"
)

var (
	vagrantMinVersion = version.Must(version.NewVersion("1.7.99"))
)

// Project returns the hashitools Project for this.
func Project(ctx *context.Shared) *hashitools.Project {
	return &hashitools.Project{
		Name:       "vagrant",
		MinVersion: vagrantMinVersion,
		Installer: &hashitools.VagrantInstaller{
			Ui: ctx.Ui,
		},
	}
}

// Vagrant wraps `vagrant` execution into an easy-to-use API.
type Vagrant struct {
	// Dir is the working directory where all Vagrant commands will
	// be executed from.
	Dir string

	// DataDir is the directory where Vagrant commands should store data.
	DataDir string

	// Env is extra environment variables to set when executing Vagrant.
	// This will be on top of the environment variables that are in this
	// process.
	Env map[string]string

	// Ui, if given, will be used to stream output from the Vagrant
	// commands. If this is nil, then the output will be logged but
	// won't be visible to the user.
	Ui ui.Ui

	// Callbacks is a mapping of callbacks that will be called for certain
	// event types within the output. These will always be serialized and
	// will block on this callback returning so it is important to make
	// this fast.
	Callbacks map[string]OutputCallback

	lock sync.Mutex
}

// A global mutex to prevent any Vagrant commands from running in parallel,
// which is not a supported mode of operation for Vagrant.
var vagrantMutex = &sync.Mutex{}

const (
	// The environment variable that Vagrant uses to configure its working dir
	vagrantCwdEnvVar = "VAGRANT_CWD"

	// The environment variable that Vagrant uses to configure its data dir.
	vagrantDataDirEnvVar = "VAGRANT_DOTFILE_PATH"
)

// Execute executes a raw Vagrant command.
func (v *Vagrant) Execute(command ...string) error {
	vagrantMutex.Lock()
	defer vagrantMutex.Unlock()

	if v.Env == nil {
		v.Env = make(map[string]string)
	}

	// Where to store data
	v.Env[vagrantDataDirEnvVar] = v.DataDir

	// Make sure we use our cwd properly
	v.Env[vagrantCwdEnvVar] = v.Dir

	// Build up the environment
	env := os.Environ()
	for k, v := range v.Env {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	// Build the command to execute
	cmd := exec.Command("vagrant", command...)
	cmd.Dir = v.Dir
	cmd.Env = env

	// Build our custom UI that we'll use that'll call the registered
	// callbacks as well as streaming data to the UI.
	callbacks := make(map[string]OutputCallback)
	callbacks["ui"] = v.uiCallback
	for n, cb := range v.Callbacks {
		callbacks[n] = cb
	}
	ui := &vagrantUi{Callbacks: callbacks}

	// Run it with the execHelper
	err := execHelper.Run(v.Ui, cmd)
	ui.Finish()
	if err != nil {
		return fmt.Errorf(
			"Error executing Vagrant: %s\n\n"+
				"The error messages from Vagrant are usually very informative.\n"+
				"Please read it carefully and fix any issues it mentions. If\n"+
				"the message isn't clear, please report this to the Otto project.",
			err)
	}

	return nil
}

func (v *Vagrant) ExecuteSilent(command ...string) error {
	v.lock.Lock()
	defer v.lock.Unlock()

	// Store the old UI and restore it before exit
	old := v.Ui
	defer func() { v.Ui = old }()

	// Make the Ui silent
	v.Ui = &ui.Logged{Ui: &ui.Null{}}
	return v.Execute(command...)
}

func (v *Vagrant) uiCallback(o *Output) {
	// If we don't have a UI return
	if v.Ui == nil {
		v.Ui = &ui.Logged{Ui: &ui.Null{}}
	}

	// Output the things to our own UI!
	v.Ui.Raw(o.Data[1] + "\n")
}
