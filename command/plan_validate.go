package command

import (
	"fmt"
	"strings"

	"github.com/hashicorp/otto/helper/flag"
	"github.com/hashicorp/otto/otto"
	"github.com/hashicorp/otto/plan"
	"github.com/mitchellh/cli"
)

// PlanValidateCommand is the command that plans any deployment.
type PlanValidateCommand struct {
	Meta
}

func (c *PlanValidateCommand) Run(args []string) int {
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

	// Parse the plan
	raw, err := plan.ParseFile(args[0])
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

	opts := &otto.PlanOpts{Validate: true}
	if err := plan.Execute(core, opts); err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	c.Ui.Output("Plan validated successfully.")
	return 0
}

func (c *PlanValidateCommand) Synopsis() string {
	return "Validate a plan created by Otto"
}

func (c *PlanValidateCommand) Help() string {
	helpText := `
Usage: otto plan validate [options] PATH

  Validate a plan created by Otto for this Appfile.

  This does basic validation on a plan created by Otto. The validation
  will check syntax, verify parameters as much as possible, and verify
  that tasks are usable.

  Passing this validation does not guarantee that plan execution will
  work successfully, but it does promise that it will begin execution
  properly.

  This can only validate plans created for this Appfile. Arbitrary plans
  from other projects can't be validated.

`

	return strings.TrimSpace(helpText)
}
