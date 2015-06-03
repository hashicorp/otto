package command

import (
	"strings"
)

// CompileCommand is the command that is responsible for "compiling" the
// Appfile into a set of data that is used by the other commands for
// execution.
type CompileCommand struct {
	Meta
}

func (c *CompileCommand) Run(args []string) int {
	return 0
}

func (c *CompileCommand) Synopsis() string {
	return "Prepares your project for being run."
}

func (c *CompileCommand) Help() string {
	helpText := `
Usage: otto [options] [path]

  Compiles the Appfile into the set of supporting files used for
  development, deploy, etc. If path is not specified, the current directory
  is assumed.

`

	return strings.TrimSpace(helpText)
}
