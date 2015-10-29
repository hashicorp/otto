package command

import (
	"os"
	"strings"

	"github.com/hashicorp/otto/plugin"
)

// PluginBuiltinCommand is a command for serving our internal plugins.
// This is not a command we expect users to call directly.
type PluginBuiltinCommand struct {
	Meta
}

func (c *PluginBuiltinCommand) Run(args []string) int {
	os.Args = append([]string{"self"}, args...)
	plugin.ServeMux(c.PluginMap)
	return 0
}

func (c *PluginBuiltinCommand) Synopsis() string {
	return "Runs a built-in plugin"
}

func (c *PluginBuiltinCommand) Help() string {
	helpText := `
Usage: otto plugin-builtin TYPE NAME

  Runs a built-in plugin.

  This is an internal command that shouldn't be called directly.

`

	return strings.TrimSpace(helpText)
}
