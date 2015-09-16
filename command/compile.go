package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/appfile/detect"
	"github.com/hashicorp/otto/ui"
)

// CompileCommand is the command that is responsible for "compiling" the
// Appfile into a set of data that is used by the other commands for
// execution.
type CompileCommand struct {
	Meta

	Detectors []*detect.Detector
}

func (c *CompileCommand) Run(args []string) int {
	var flagAppfile string
	fs := c.FlagSet("compile", FlagSetNone)
	fs.Usage = func() { c.Ui.Error(c.Help()) }
	fs.StringVar(&flagAppfile, "appfile", "", "")
	if err := fs.Parse(args); err != nil {
		return 1
	}

	// Load a UI
	ui := c.OttoUi()
	ui.Header("Loading Appfile...")

	// Determine all the Appfile paths
	//
	// First, if an Appfile was specified on the command-line, it must
	// exist so we validate that it exists.
	if flagAppfile != "" {
		fi, err := os.Stat(flagAppfile)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error loading Appfile: %s", err))
			return 1
		}

		if fi.IsDir() {
			flagAppfile = filepath.Join(flagAppfile, DefaultAppfile)
		}
	}

	// If the Appfile is still blank, just use our current directory
	if flagAppfile == "" {
		var err error
		flagAppfile, err = os.Getwd()
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error loading working directory: %s", err))
			return 1
		}

		flagAppfile = filepath.Join(flagAppfile, DefaultAppfile)
	}

	// If we have the Appfile, then make sure it is an absoute path
	if flagAppfile != "" {
		var err error
		flagAppfile, err = filepath.Abs(flagAppfile)
		if err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error getting Appfile path: %s", err))
			return 1
		}
	}

	// Load the appfile. This is the only time we ever load the
	// raw Appfile. All other commands load the compiled Appfile.
	var app *appfile.File
	if fi, err := os.Stat(flagAppfile); err == nil && !fi.IsDir() {
		app, err = appfile.ParseFile(flagAppfile)
		if err != nil {
			c.Ui.Error(err.Error())
			return 1
		}
	}

	// Tell the user what is happening if they have no Appfile
	if app == nil {
		ui.Header("No Appfile found! Detecting project information...")
		ui.Message(fmt.Sprintf(
			"No Appfile was found. If there is no Appfile, Otto will do its best\n" +
				"to detect the type of application this is and set reasonable defaults.\n" +
				"This is a good way to get started with Otto, but over time we recommend\n" +
				"writing a real Appfile since this will allow more complex customizations,\n" +
				"the ability to reference dependencies, versioning, and more."))
	}

	// Load the default Appfile so we can merge in any defaults into
	// the loaded Appfile (if there is one).
	detectConfig := &detect.Config{
		Detectors: c.Detectors,
	}
	appDef, err := appfile.Default(filepath.Dir(flagAppfile), detectConfig)
	if err != nil {
		c.Ui.Error(fmt.Sprintf(
			"Error loading Appfile: %s", err))
		return 1
	}

	// If there was no loaded Appfile and we don't have an application
	// type then we weren't able to detect the type. Error.
	if app == nil && appDef.Application.Type == "" {
		c.Ui.Error(strings.TrimSpace(errCantDetectType))
		return 1
	}

	// Merge the appfiles
	if app != nil {
		if err := appDef.Merge(app); err != nil {
			c.Ui.Error(fmt.Sprintf(
				"Error loading Appfile: %s", err))
			return 1
		}
	}
	app = appDef

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
		case *appfile.CompileEventImport:
			ui.Message(fmt.Sprintf(
				"Fetching import: %s", e.Source))
		}
	}
}

const errCantDetectType = `
No Appfile is present and Otto couldn't detect the project type automatically.
Otto does its best without an Appfile to detect what kind of project this is
automatically, but sometimes this fails if the project is in a structure
Otto doesn't recognize or its a project type that Otto doesn't yet support.

Please create an Appfile and specify at a minimum the project type. Below
is an example minimal Appfile specifying the "go" project type:

    application {
        type = "go"
    }

If you believe Otto should've been able to automatically detect your
project type, then please open an issue with the Otto project.
`
