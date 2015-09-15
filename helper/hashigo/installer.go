package hashigo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/otto/ui"
	"gopkg.in/flosch/pongo2.v3"
)

// Installer is the interface that knows how to install things.
//
// This is an interface to support different installation methods between
// our different projects.
type Installer interface {
	// InstallAsk should ask the user if they'd like to install the
	// project. This is only called if installation is actually required.
	InstallAsk(installed, required, latest *version.Version) (bool, error)

	// Install should install the specified version.
	Install(*version.Version) error

	// Path is the path to the installed main binary of this project,
	// or "" if it doesn't seem installed.
	Path() string
}

// GoInstaller is an Installer that knows how to install Go projects.
type GoInstaller struct {
	// Name is the name of the project to install
	Name string

	// Dir is the directory where projects will be installed. They will
	// be installed to a sub-directory of the project name. Example:
	// if Dir is "/foo", then the Packer binary would be installed to
	// "/foo/packer/packer"
	Dir string

	// Ui is the Otto UI for asking the user for input and outputting
	// the status of installation.
	Ui ui.Ui
}

func (i *GoInstaller) InstallAsk(installed, required, latest *version.Version) (bool, error) {
	input := &ui.InputOpts{
		Id:      fmt.Sprintf("%s_install", i.Name),
		Query:   fmt.Sprintf("Would you like Otto to install %s?", strings.Title(i.Name)),
		Default: "",
	}

	// Figure out the description text to use for input
	var tplString string
	if installed == nil {
		tplString = installRequired
	} else {
		tplString = installRequiredUpdate
	}

	// Parse the template and render it
	tpl, err := pongo2.FromString(strings.TrimSpace(tplString))
	if err != nil {
		return false, err
	}
	input.Description, err = tpl.Execute(map[string]interface{}{
		"name":      i.Name,
		"installed": installed.String(),
		"latest":    latest.String(),
		"required":  required.String(),
	})
	if err != nil {
		return false, err
	}

	result, err := i.Ui.Input(input)
	if err != nil {
		return false, err
	}

	return strings.ToLower(result) == "yes", nil
}

func (i *GoInstaller) Install(vsn *version.Version) error {
	return fmt.Errorf("error")
}

func (i *GoInstaller) Path() string {
	path := filepath.Join(i.Dir, i.Name, i.Name)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	return ""
}

const installRequired = `
Otto requires {{name}} to be installed, but it couldn't be found on your
system. Otto can install the latest version of {{name}} for you. Otto will
install this into its own private data directory so it doesn't conflict
with anything else on your system. Would you like Otto to install {{name}}
for you? Alternatively, you may install this on your own.

Please enter 'yes' to continue. Any other value will exit.
`

const installRequiredUpdate = `
An older version of {{name}} was found installed ({{installed}}). Otto requires
version {{required}} or higher. Otto can install the latest version of {{name}}
for you ({{latest}}). Would you like Otto to install {{name}} for you?

Please enter 'yes' to continue. Any other value will exit.
`
