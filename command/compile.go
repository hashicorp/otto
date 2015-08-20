package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/ui"
)

// CompileCommand is the command that is responsible for "compiling" the
// Appfile into a set of data that is used by the other commands for
// execution.
type CompileCommand struct {
	Meta
}

func (c *CompileCommand) Run(args []string) int {
	var flagAppfile string
	fs := c.FlagSet("compile", FlagSetNone)
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	fs.StringVar(&flagAppfile, "appfile", os.Getenv(EnvAppFile), "")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Load a UI
	ui := c.OttoUi()

	// Load the appfile. This is the only time we ever load the
	// raw Appfile. All other commands load the compiled Appfile.
	if flagAppfile == "" {
		flagAppfile = "."
	}
	fi, err := os.Stat(flagAppfile)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error checking Appfile path: %s", err))
		return 1
	}
	if fi.IsDir() {
		flagAppfile = filepath.Join(flagAppfile, DefaultAppfile)
	}
	app, err := appfile.ParseFile(flagAppfile)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	// Compile the Appfile
	ui.Header("Fetching all Appfile dependencies...")
	capp, err := appfile.Compile(app, &appfile.CompileOpts{
		Dir: filepath.Join(
			filepath.Dir(app.Path), DefaultOutputDir, DefaultOutputDirCompiledAppfile),
		Callback: c.compileCallback(ui),
	})
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error compiling Appfile: %s", err))
		return 1
	}

	// Get a core
	core, err := c.Core(capp)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error loading core: %s", err))
		return 1
	}

	// Get the active infrastructure just for UI reasons
	infra := app.ActiveInfrastructure()

	// Before the compilation, output to the user what is going on
	ui.Header("Compiling...")
	ui.Message(fmt.Sprintf(
		"Application:    %s (%s)",
		app.Application.Name,
		app.Application.Type))
	ui.Message(fmt.Sprintf("Project:        %s", app.Project.Name))
	ui.Message(fmt.Sprintf(
		"Infrastructure: %s (%s)",
		infra.Type,
		infra.Flavor))
	ui.Message("")

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

func (c *CompileCommand) compileCallback(ui ui.Ui) func(appfile.CompileEvent) {
	return func(raw appfile.CompileEvent) {
		switch e := raw.(type) {
		case *appfile.CompileEventDep:
			ui.Message(fmt.Sprintf(
				"Fetching dependency: %s", e.Source))
		}
	}
}
