package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/otto/helper/flag"
	"github.com/hashicorp/otto/otto"
	"github.com/hashicorp/otto/plan"
	"github.com/mitchellh/cli"
)

// PlanExecuteCommand is the command that plans any deployment.
type PlanExecuteCommand struct {
	Meta
}

func (c *PlanExecuteCommand) Run(args []string) int {
	fs := c.FlagSet("plan", FlagSetNone)
	incArgs, _, args := flag.FilterArgs(fs, args)
	if err := fs.Parse(incArgs); err != nil {
		return cli.RunResultHelp
	}

	// Verify we got a path
	if len(args) != 1 {
		c.Ui.Error(fmt.Sprintf(
			"Expected a single argument: path"))
		return cli.RunResultHelp
	}

	// Determine if we're getting the plan from stdin or a file
	var r io.Reader = os.Stdin
	if args[0] != "-" {
		f, err := os.Open(args[0])
		if err != nil {
			c.Ui.Error(err.Error())
			return 1
		}

		r = f
	}

	// Parse the plan
	raw, err := plan.Parse(r)
	if c, ok := r.(io.Closer); ok {
		c.Close()
	}
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	plan := otto.Plan{Plans: raw}

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

	opts := &otto.PlanOpts{Validate: false}
	if err := plan.Execute(core, opts); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	return 0
}

func (c *PlanExecuteCommand) Synopsis() string {
	return "Execute a plan created by Otto"
}

func (c *PlanExecuteCommand) Help() string {
	helpText := `
Usage: otto plan execute [options] PATH

  Execute a plan created by Otto for this Appfile.

  This can only execute plans created for this Appfile. Arbitrary plans
  from other projects may result in unexpected behavior.

  If the PATH argument is "-" then this will read from stdin.

`

	return strings.TrimSpace(helpText)
}
