package command

import (
	"fmt"
	"strings"
)

// StatusCommand is the command that shows the status of the various
// stages of this application.
type StatusCommand struct {
	Meta
}

func (c *StatusCommand) Run(args []string) int {
	fs := c.FlagSet("status", FlagSetNone)
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

	// Execute the task
	err = core.Status()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error occurred: %s", err))
		return 1
	}

	return 0
}

func (c *StatusCommand) Synopsis() string {
	return "Status of the stages of this application"
}

func (c *StatusCommand) Help() string {
	helpText := `
Usage: otto status

  Shows the status of this application.

  This command will show whether an application has a development
  environment created, a build made, a deploy done, the infrastructure
  ready, etc.

  The output from this command is loaded from a cache and may not represent
  the true state of the world if external changes have happened outside of
  Otto. For a true status, the actual command to manage each thing must
  be run. For example, for development "otto dev" should be run.

`

	return strings.TrimSpace(helpText)
}
