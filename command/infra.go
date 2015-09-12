package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/otto/helper/flag"
)

// InfraCommand is the command that sets up the infrastructure for an
// Appfile.
type InfraCommand struct {
	Meta
}

func (c *InfraCommand) Run(args []string) int {
	fs := c.FlagSet("infra", FlagSetNone)
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	args, execArgs, posArgs := flag.FilterArgs(fs, args)
	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Get the remaining args to determine if we have an action.
	var action string
	if len(posArgs) > 0 {
		action = posArgs[0]
		execArgs = append(execArgs, posArgs[1:]...)
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
	err = core.Infra(action, execArgs)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error occurred: %s", err))
		return 1
	}

	return 0
}

func (c *InfraCommand) Synopsis() string {
	return "Builds the infrastructure for the Appfile"
}

func (c *InfraCommand) Help() string {
	helpText := `
Usage: otto infra [options]

  Builds the infrastructure for the Appfile.

  This will create real infrastructure resources as configured by the
  Appfile, such as launching real servers. This command is stateful. If
  the infrastructure has already been created, it won't create it again.
  If the infrastructure is created but needs to be modified, it will be
  modified.

  Note that not all infrastructure changes are non-destructive and this
  command may cause downtime.

`

	return strings.TrimSpace(helpText)
}
