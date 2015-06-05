package command

import (
	"fmt"
	"strings"
)

// CompileCommand is the command that is responsible for "compiling" the
// Appfile into a set of data that is used by the other commands for
// execution.
type CompileCommand struct {
	Meta
}

func (c *CompileCommand) Run(args []string) int {
	fs := c.FlagSet("compile", FlagSetAppfile|FlagSetOutputDir)
	if err := fs.Parse(args); err != nil {
		return 1
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

	// Get the active infrastructure just for UI reasons
	infra := app.ActiveInfrastructure()

	// Before the compilation, output to the user what is going on
	ui := c.OttoUi()
	ui.Header("Compiling...")
	ui.Message(fmt.Sprintf("Application:    %s", app.Application.Name))
	ui.Message(fmt.Sprintf("Project:        %s", app.Project.Name))
	ui.Message(fmt.Sprintf(
		"Infrastructure: %s (%s)",
		infra.Type,
		infra.Flavor))

	// Compile!
	if err := core.Compile(); err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error compiling: %s", err))
		return 1
	}

	// Success!
	ui.Header("Compilation success!")
	ui.Message(fmt.Sprintf(
		"This means that Otto is now ready to start a development environment,\n" +
			"deploy this application, build the supporting infastructure, and\n" +
			"more. See the help for more information.\n\n" +
			"Supporting files to enable Otto to manage your application from\n" +
			"development to deployment have been placed in the output directory.\n" +
			"These files can be manually inspected to determine what Otto will do."))

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
