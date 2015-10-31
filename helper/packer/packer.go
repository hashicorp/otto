package packer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/otto/context"
	execHelper "github.com/hashicorp/otto/helper/exec"
	"github.com/hashicorp/otto/helper/hashitools"
	"github.com/hashicorp/otto/ui"
)

var (
	packerMinVersion = version.Must(version.NewVersion("0.8.6"))
)

// Project returns the hashitools Project for this.
func Project(ctx *context.Shared) *hashitools.Project {
	return &hashitools.Project{
		Name:       "packer",
		MinVersion: packerMinVersion,
		Installer: &hashitools.GoInstaller{
			Name: "packer",
			Dir:  filepath.Join(ctx.InstallDir),
			Ui:   ctx.Ui,
		},
	}
}

// Packer wraps `packer` execution into an easy-to-use API
type Packer struct {
	// Path is the path to Packer itself. If empty, "packer"
	// will be used and looked up via the PATH var.
	Path string

	// Dir is the working directory where all Packer commands are executed
	Dir string

	// Ui, if given, will be used to stream output from the Packer commands.
	// If this is nil, then the output will be logged but won't be visible
	// to the user.
	Ui ui.Ui

	// Callbacks is a list of callbacks that will be called for certain
	// event types within the output
	Callbacks map[string]OutputCallback

	// Variables is a list of variables to pass to Packer.
	Variables map[string]string
}

// Execute executes a raw Packer command.
func (p *Packer) Execute(commandRaw ...string) error {
	varfile, err := p.varfile()
	if err != nil {
		return err
	}
	if execHelper.ShouldCleanup() {
		defer os.Remove(varfile)
	}

	// The command must always be machine-readable. We use this
	// exclusively to mirror the UI output.
	command := make([]string, len(commandRaw)+3)
	command[0] = commandRaw[0]
	command[1] = "-machine-readable"
	command[2] = "-var-file"
	command[3] = varfile
	copy(command[4:], commandRaw[1:])

	// Build the command to execute
	path := "packer"
	if p.Path != "" {
		path = p.Path
	}
	cmd := exec.Command(path, command...)
	cmd.Dir = p.Dir

	// Build our custom UI that we'll use that'll call the registered
	// callbacks as well as streaming data to the UI.
	callbacks := make(map[string]OutputCallback)
	callbacks["ui"] = p.uiCallback
	for n, cb := range p.Callbacks {
		callbacks[n] = cb
	}
	ui := &packerUi{Callbacks: callbacks}

	// Execute!
	err = execHelper.Run(ui, cmd)
	ui.Finish()
	if err != nil {
		return fmt.Errorf(
			"Error executing Packer: %s\n\n"+
				"The error messages from Packer are usually very informative.\n"+
				"Please read it carefully and fix any issues it mentions. If\n"+
				"the message isn't clear, please report this to the Otto project.",
			err)
	}

	return nil
}

func (p *Packer) uiCallback(o *Output) {
	// If we don't have a UI return
	if p.Ui == nil {
		p.Ui = &ui.Logged{Ui: &ui.Null{}}
	}

	// Output the things to our own UI!
	p.Ui.Raw(o.Data[1] + "\n")
}

func (p *Packer) varfile() (string, error) {
	f, err := ioutil.TempFile("", "otto")
	if err != nil {
		return "", err
	}

	err = json.NewEncoder(f).Encode(p.Variables)
	f.Close()
	if err != nil {
		os.Remove(f.Name())
	}

	return f.Name(), err
}
