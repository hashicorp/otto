package command

import (
	"strings"
)

// InfraCommand is the command that sets up the infrastructure for an
// Appfile.
type InfraCommand struct {
	Meta
}

func (c *InfraCommand) Run(args []string) int {
	fs := c.FlagSet("infra", FlagSetAppfile)
	if err := fs.Parse(args); err != nil {
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
