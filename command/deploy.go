package command

import (
	"fmt"
	"strings"
)

// DeployCommand is the command that deploys the app once it is built.
type DeployCommand struct {
	Meta
}

func (c *DeployCommand) Run(args []string) int {
	fs := c.FlagSet("deploy", FlagSetNone)
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

	// Deploy the artifact
	if err := core.Deploy(); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error building app: %s", err))
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
