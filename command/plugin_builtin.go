package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/otto/builtin/pluginmap"
	"github.com/hashicorp/otto/plugin"
)

// PluginBuiltinCommand is a command for serving our internal plugins.
// This is not a command we expect users to call directly.
type PluginBuiltinCommand struct {
	Meta
}

func (c *PluginBuiltinCommand) Run(args []string) int {
	fs := c.FlagSet("plugin-builtin", FlagSetNone)
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	if len(args) != 2 {
		fs.Usage()
		return 1
	}

	var opts plugin.ServeOpts
	switch args[0] {
	case "app":
		opts.AppFunc = pluginmap.Apps[args[1]]
	default:
		c.Ui.Error(fmt.Sprintf("Unknown plugin type: %s", args[0]))
		return 1
	}

	plugin.Serve(&opts)
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
