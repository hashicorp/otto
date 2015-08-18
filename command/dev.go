package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/otto/helper/flag"
	"github.com/hashicorp/otto/otto"
)

// DevCommand is the command that manages (starts, stops, etc.) the
// development environment for an Appfile.
type DevCommand struct {
	Meta
}

func (c *DevCommand) Run(args []string) int {
	fs := c.FlagSet("dev", FlagSetNone)
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
	err = core.Execute(&otto.ExecuteOpts{
		Task:   otto.ExecuteTaskDev,
		Action: action,
		Args:   execArgs,
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error occurred: %s", err))
		return 1
	}

	return 0
}

func (c *DevCommand) Synopsis() string {
	return "Start and manage a development environment"
}

func (c *DevCommand) Help() string {
	helpText := `
Usage: otto dev [options]

  Start and manage a development environment for your application.

  This will start a development environment for your application.
  Additional subcommands such as "destroy" can be given to tear down
  the development environment.

  The development environment will be local and will automatically include
  all upstream dependencies within the environment properly configured
  and started.

`

	return strings.TrimSpace(helpText)
}
