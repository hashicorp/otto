package command

import (
	"strings"
)

// HelpCommand is not a real command. It just shows help output for
// people expecting `otto help` to work.
type HelpCommand struct {
	Meta
}

func (c *HelpCommand) Run(args []string) int {
	c.Ui.Error(strings.TrimSpace(helpOutput))
	return 1
}

func (c *HelpCommand) Synopsis() string {
	return "Not a real command."
}

func (c *HelpCommand) Help() string {
	return "Not a real command."
}

const helpOutput = `
Otto doesn't use "otto help" for subcommand help!

For a list of all available top level subcommands, run "otto -h".

For individual commands use the "-h" flag to get help. For example:
"otto status -h" will show you how to use the status command.

For commands that take subcommands such as "dev", use the "help" subcommand
to get a full listing of availabile subcommands. For example: "otto dev help".
`
