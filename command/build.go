package command

import (
	"fmt"
	"strings"
)

// BuildCommand is the command that builds a deployable artifact
// for this version of the app.
type BuildCommand struct {
	Meta
}

func (c *BuildCommand) Run(args []string) int {
	fs := c.FlagSet("build", FlagSetNone)
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Load the appfile
	app, err := c.Appfile()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	// Get a core
	core, err := c.Core(app)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error loading core: %s", err))
		return 1
	}

	// Build the artifact
	if err := core.Build(); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error building app: %s", err))
		return 1
	}

	return 0
}

func (c *BuildCommand) Synopsis() string {
	return "Build the deployable artifact for the app"
}

func (c *BuildCommand) Help() string {
	helpText := `
Usage: otto build [options]

  Builds the deployable artifact for the app on the target
  infrastructure specified during compilation of the Appfile.

  This will build and inventory the artifact that is deployable
  for the app represented by this Appfile.

`

	return strings.TrimSpace(helpText)
}
