package terraform

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	execHelper "github.com/hashicorp/otto/helper/exec"
	"github.com/hashicorp/otto/ui"
)

// Terraform wraps `terraform` execution into an easy-to-use API
type Terraform struct {
	// Dir is the working directory where all Terraform commands are executed
	Dir string

	// Ui, if given, will be used to stream output from the Terraform commands.
	// If this is nil, then the output will be logged but won't be visible
	// to the user.
	Ui ui.Ui

	// Variables is a list of variables to pass to Terraform.
	Variables map[string]string
}

// Execute executes a raw Terraform command
func (t *Terraform) Execute(commandRaw ...string) error {
	command := commandRaw

	// If we have variables, create the var file
	if len(t.Variables) > 0 {
		varfile, err := t.varfile()
		if err != nil {
			return err
		}
		defer os.Remove(varfile)

		// Append the varfile onto our command.
		command = append(command, "-var-file", varfile)
	}

	// Build the command to execute
	cmd := exec.Command("terraform", command...)
	cmd.Dir = t.Dir

	// Start the Terraform command
	if err := execHelper.Run(t.Ui, cmd); err != nil {
		err = fmt.Errorf("Error running Terraform: %s", err)
		return err
	}

	return nil
}

// Outputs reads the outputs from the given state path. It returns
// nil if there are no outputs.
func (t *Terraform) Outputs(path string) (map[string]string, error) {
	if _, err := os.Stat(path); err != nil {
		// If the file doesn't exist, we just don't have outputs
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	return Outputs(path)
}

func (t *Terraform) varfile() (string, error) {
	f, err := ioutil.TempFile("", "otto-tf")
	if err != nil {
		return "", err
	}

	err = json.NewEncoder(f).Encode(t.Variables)
	f.Close()
	if err != nil {
		os.Remove(f.Name())
	}

	return f.Name(), err
}
