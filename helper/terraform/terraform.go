package terraform

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/otto/directory"
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

	// Directory can be set to point to a directory where data can be
	// stored. If this is set, then the state will be loaded/stored here
	// automatically.
	//
	// StateId is the identifier used to load/store the state from
	// blob storage. If this is empty, state won't be loaded or stored
	// automatically.
	//
	// It is highly recommended to use this instead of manually attempting
	// to manage state since this will properly handle storing the state
	// in the issue of an error and will put the state in the pwd in the
	// case we can't write it to a directory.
	Directory directory.Backend
	StateId   string
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

	// If we care about state, then setup the state directory and
	// load it up.
	var stateDir string
	statePath := filepath.Join(stateDir, "state")
	if t.StateId != "" && t.Directory != nil {
		var err error
		stateDir, err = ioutil.TempDir("", "otto-tf")
		if err != nil {
			return err
		}
		defer os.RemoveAll(stateDir)

		// State path
		stateOldPath := filepath.Join(stateDir, "state.old")

		// Load the state from the directory
		data, err := t.Directory.GetBlob(t.StateId)
		if err != nil {
			return fmt.Errorf("Error loading Terraform state: %s", err)
		}
		if data != nil {
			err = data.WriteToFile(stateOldPath)
			data.Close()
		}
		if err != nil {
			return fmt.Errorf("Error writing Terraform state: %s", err)
		}

		// Append the state to the args
		command = append(command, "-state", stateOldPath)
		command = append(command, "-state-out", statePath)
	}

	// Build the command to execute
	cmd := exec.Command("terraform", command...)
	cmd.Dir = t.Dir

	// Start the Terraform command. If there is an error we just store
	// the error but can't exit yet because we have to store partial
	// state if there is any.
	err := execHelper.Run(t.Ui, cmd)
	if err != nil {
		err = fmt.Errorf("Error running Terraform: %s", err)
	}

	// Save the state file if we have it.
	if f, ferr := os.Open(statePath); ferr == nil {
		// Store the state
		derr := t.Directory.PutBlob(t.StateId, &directory.BlobData{
			Data: f,
		})

		// Always close the file
		f.Close()

		// If we couldn't save the data, then note the error. This is a
		// _really_ bad error to get since there isn't a good way to
		// recover. For now, we just copy the state to the pwd and note
		// the user.
		if derr != nil {
			// TODO: copy state

			err = fmt.Errorf(
				"Failed to save Terraform state: %s\n\n",
				"This means that Otto was unable to store the state of your infrastructure.\n"+
					"At this time, Otto doesn't support gracefully recovering from this\n"+
					"scenario. The state should be in the path below. Please ask the\n"+
					"community for assistance.",
				derr)
		}
	}

	return err
}

// Outputs reads the outputs from the configured directory storage.
func (t *Terraform) Outputs() (map[string]string, error) {
	// Make a temporary file to store our state
	tf, err := ioutil.TempFile("", "otto-tf")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tf.Name())

	// Read the state from the directory and put it on disk. Lots of
	// careful management of file handles here.
	data, err := t.Directory.GetBlob(t.StateId)
	if err == nil {
		if data == nil {
			return nil, nil
		}

		_, err = io.Copy(tf, data.Data)
		data.Close()
	}
	tf.Close()
	if err != nil {
		return nil, fmt.Errorf("Error loading Terraform state: %s", err)
	}

	// Read the outputs as normal. Defers will clean up our temp file.
	return Outputs(tf.Name())
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
