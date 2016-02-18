package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/otto/helper/flag"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
)

// InfrasCommand is the command that lists and shows apps in the directory.
type InfrasCommand struct {
	Meta
}

func (c *InfrasCommand) Run(args []string) int {
	fs := c.FlagSet("infras", FlagSetNone)
	args, _, _ = flag.FilterArgs(fs, args)
	if err := fs.Parse(args); err != nil {
		return cli.RunResultHelp
	}

	/*
		// TODO: Get the remaining args to determine if we have an ID to view
		if len(posArgs) > 0 {
			action = posArgs[0]
			execArgs = append(execArgs, posArgs[1:]...)
		}
	*/

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

	// List the infras
	infras, err := core.Directory().ListInfra()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if len(infras) == 0 {
		c.Ui.Output(strings.TrimSpace(outputNoInfras))
		return 1
	}

	// Output all the infras
	output := make([]string, 1, len(infras)+1)
	output[0] = "Name | Type | Flavor | Status"
	for _, v := range infras {
		output = append(output, fmt.Sprintf(
			"%s | %s | %s | %s", v.Name, v.Type, v.Flavor, v.State))
	}

	c.Ui.Output(columnize.SimpleFormat(output))
	return 0
}

func (c *InfrasCommand) Synopsis() string {
	return "View the infrastructure managed by Otto"
}

func (c *InfrasCommand) Help() string {
	helpText := `
Usage: otto infras [id]

  List all the infrastructures or view a single infrastructure
  managed by Otto.

  Otto maintains metadata about the infrastructures that were built or
  are managed by Otto. This command can list and inspect these
  infrastructures.

  By default Otto stores all data in a single directory backend global to
  your logged in operating system user. Please see the online documentation
  on directory backends for how to configure more. Directory backends should
  be configured for your organization.

`

	return strings.TrimSpace(helpText)
}

const outputNoInfras = `
There are currently no infrastructures managed by Otto!

You can create an infrastructure by configuring an Appfile, compiling it
with "otto compile", and then deploying that application with "otto deploy".
If you only want to deploy the infrastructure, see the help on "otto deploy".
`
