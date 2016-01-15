package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/otto/helper/flag"
	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
)

// AppsCommand is the command that lists and shows apps in the directory.
type AppsCommand struct {
	Meta
}

func (c *AppsCommand) Run(args []string) int {
	fs := c.FlagSet("apps", FlagSetNone)
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

	// List the applications
	apps, err := core.Directory().ListApps()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	if len(apps) == 0 {
		c.Ui.Output("FOO")
		return 1
	}

	// Output all the apps
	output := make([]string, 1, len(apps)+1)
	output[0] = "Name | Version | ID"
	for _, a := range apps {
		output = append(output, fmt.Sprintf(
			"%s | %s | %s", a.Name, a.AppLookup.Version, a.AppLookup.AppID))
	}

	c.Ui.Output(columnize.SimpleFormat(output))
	return 0
}

func (c *AppsCommand) Synopsis() string {
	return "View the apps managed by Otto"
}

func (c *AppsCommand) Help() string {
	helpText := `
Usage: otto apps [id]

  List all the applications or view a single application managed by Otto.

  Otto maintains metadata about the applications that were compiled
  for this directory. This command can list and inspect these applications.

  By default Otto stores all data in a single directory backend global to
  your logged in operating system user. Please see the online documentation
  on directory backends for how to configure more. Directory backends should
  be configured for your organization.

`

	return strings.TrimSpace(helpText)
}
