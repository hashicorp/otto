package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/otto/helper/flag"
)

// DeployCommand is the command that deploys the app once it is built.
type DeployCommand struct {
	Meta
}

func (c *DeployCommand) Run(args []string) int {
	fs := c.FlagSet("deploy", FlagSetNone)
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

	// Destroy action gets an extra double-check
	if action == "destroy" {
		msg := "Otto will delete all resources associated with the deploy."
		if !c.confirmDestroy(msg, execArgs) {
			return 1
		}
	}

	// Deploy the artifact
	if err := core.Deploy(action, execArgs); err != nil {
		// Display errors without prefix, we expect them to be formatted in a way
		// that's suitable for UI.
		c.Ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *DeployCommand) Synopsis() string {
	return "Deploy the application"
}

func (c *DeployCommand) Help() string {
	helpText := `
Usage: otto deploy [options]

  Deploy the application to the current environment.

  This command may modify real infrastructure, including adding
  and removing resources.

  The "build" command should be called prior to this to create the
  build artifact. Deploy can be called multiple times with the same
  artifact to redeploy an application.

`

	return strings.TrimSpace(helpText)
}
