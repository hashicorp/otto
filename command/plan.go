package command

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/otto/helper/flag"
	"github.com/hashicorp/otto/plan"
	"github.com/mitchellh/cli"
)

// PlanCommand is the command that plans any deployment.
type PlanCommand struct {
	Meta
}

func (c *PlanCommand) Run(args []string) int {
	var flagOut string
	fs := c.FlagSet("plan", FlagSetNone)
	fs.StringVar(&flagOut, "out", "", "")
	args, _, _ = flag.FilterArgs(fs, args)
	if err := fs.Parse(args); err != nil {
		return cli.RunResultHelp
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

	// Get the plan
	p, err := core.Plan()
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error occurred: %s", err))
		return 1
	}

	// No matter what we get, output a plan if we request a file
	if flagOut != "" {
		var w io.Writer = os.Stdout
		if flagOut != "-" {
			w, err = os.Create(flagOut)
			if err != nil {
				c.Ui.Error(fmt.Sprintf(
					"Error opening file for writing: %s", err))
				return 1
			}
		}

		err = printer.Fprint(w, plan.EncodeHCL(p.Plans))
		w.Write([]byte("\n"))
		if c, ok := w.(io.Closer); ok {
			c.Close()
		}
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error writing plan: %s", err))
			return 1
		}
	}

	// If we don't have a plan, then let the user know
	if p.Empty() {
		c.Ui.Output("Everything is up-to-date. No changes needed!")
		return 0
	}

	// Output the plan
	if len(p.Plans) > 0 {
		c.Ui.Output("Infrastructure:\n")
		for _, p := range p.Plans {
			c.Ui.Output(fmt.Sprintf("  Plan: %s", p.Description))
			for _, t := range p.Tasks {
				c.Ui.Output(fmt.Sprintf("    Task: %s", t.Type))
				c.Ui.Output(fmt.Sprintf("      Desc: %s", t.Description))
			}
		}
	}

	return 0
}

func (c *PlanCommand) Synopsis() string {
	return "Create a deployment plan"
}

func (c *PlanCommand) Help() string {
	helpText := `
Usage: otto plan [options]

  Create a deployment plan.

  This will output the plan for deploying this application. The plan
  will include any necessary infrastructure changes in addition to steps
  to simply deploy the application.

  This command will not modify any real infrastructure. This will only output
  a plan of what Otto will do. You can feed this plan directly into Otto
  to ensure that Otto only executes what is included in this plan.

  Use "otto plan execute" to execute a saved plan. To save a plan, use
  the "-out" flag.

Options:

  -out=path    Path to save the plan for manual execution later.

`

	return strings.TrimSpace(helpText)
}
