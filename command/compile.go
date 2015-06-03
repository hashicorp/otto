package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/appfile"
)

// CompileCommand is the command that is responsible for "compiling" the
// Appfile into a set of data that is used by the other commands for
// execution.
type CompileCommand struct {
	Meta
}

func (c *CompileCommand) Run(args []string) int {
	fs := c.FlagSet("compile", FlagSetNone)
	if err := fs.Parse(args); err != nil {
		return 1
	}

	args = fs.Args()

	// Get the path to where the Appfile lives
	path := "."
	if len(args) >= 1 {
		path = args[0]
	}

	// Verify the path is valid
	fi, err := os.Stat(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error checking Appfile path: %s", err))
		return 1
	}
	if fi.IsDir() {
		path = filepath.Join(path, DefaultAppfile)
	}

	// Load the appfile
	app, err := appfile.ParseFile(path)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error parsing Appfile: %s", err))
		return 1
	}

	// Get a core
	core, err := c.Core(app)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error loading core: %s", err))
		return 1
	}

	// Compile!
	if err := core.Compile(); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error compiling: %s", err))
		return 1
	}

	return 0
}

func (c *CompileCommand) Synopsis() string {
	return "Prepares your project for being run."
}

func (c *CompileCommand) Help() string {
	helpText := `
Usage: otto [options] [path]

  Compiles the Appfile into the set of supporting files used for
  development, deploy, etc. If path is not specified, the current directory
  is assumed.

`

	return strings.TrimSpace(helpText)
}
