package command

import (
	"fmt"
	"log"
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

	app, appPath, err := loadAppfile(flagAppfile)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
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

	// Parse the detectors
	dataDir, err := c.DataDir()
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	detectorDir := filepath.Join(dataDir, DefaultLocalDataDetectorDir)
	log.Printf("[DEBUG] loading detectors from: %s", detectorDir)
	detectConfig, err := detect.ParseDir(detectorDir)
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}
	if detectConfig == nil {
		detectConfig = &detect.Config{}
	}
	err = detectConfig.Merge(&detect.Config{Detectors: c.Detectors})
	if err != nil {
		c.Ui.Error(err.Error())
		return 1
	}

	// Load the default Appfile so we can merge in any defaults into
	// the loaded Appfile (if there is one).
	appDef, err := appfile.Default(appPath, detectConfig)
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
		Detect:   detectConfig,
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
	ui.Header("[green]Compilation success!")
	ui.Message(fmt.Sprintf(
		"[green]This means that Otto is now ready to start a development environment,\n" +
			"deploy this application, build the supporting infrastructure, and\n" +
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

  This command will download and update any dependencies as well as
  the import statements in your Appfile. This process only happens during
  compilation so that every other Otto operation begins executing much
  more quickly.

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

// Returns a loaded copy of any appfile.File we find, otherwise returns nil,
// which is valid, since Otto can detect app type and calculate defaults.
// Also returns the base dir of the appfile, which is the current WD in the
// case of a nil appfile.
func loadAppfile(flagAppfile string) (*appfile.File, string, error) {
	appfilePath, err := findAppfile(flagAppfile)
	if err != nil {
		return nil, "", err
	}
	if appfilePath == "" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, "", err
		}
		return nil, wd, nil
	}
	app, err := appfile.ParseFile(appfilePath)
	if err != nil {
		return nil, "", err
	}
	return app, filepath.Dir(app.Path), nil
}

// findAppfile returns the path to an existing Appfile by checking the optional
// flag value and the current directory. It returns blank if it does not find
// any Appfiles
func findAppfile(flag string) (string, error) {
	// First, if an Appfile was specified on the command-line, it must
	// exist so we validate that it exists.
	if flag != "" {
		fi, err := os.Stat(flag)
		if err != nil {
			return "", fmt.Errorf("Error loading Appfile: %s", err)
		}

		if fi.IsDir() {
			return findAppfileInDir(flag), nil
		} else {
			return flag, nil
		}
	}

	// Otherwise we search through our current directory
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Error loading working directory: %s", err)
	}
	return findAppfileInDir(wd), nil
}

// findAppfileInDir takes a path to a directory returns the path to the first
// existing Appfile by first looking for the DefaultAppfile and then looking
// for any AltAppfiles in the dir
func findAppfileInDir(path string) string {
	if _, err := os.Stat(filepath.Join(path, DefaultAppfile)); err == nil {
		return filepath.Join(path, DefaultAppfile)
	}
	for _, aaf := range AltAppfiles {
		if _, err := os.Stat(filepath.Join(path, aaf)); err == nil {
			return filepath.Join(path, aaf)
		}
	}
	return ""
}

const errCantDetectType = `
No Appfile is present and Otto couldn't detect the project type automatically.
Otto does its best without an Appfile to detect what kind of project this is
automatically, but sometimes this fails if the project is in a structure
Otto doesn't recognize or its a project type that Otto doesn't yet support.

Please create an Appfile and specify at a minimum the project name and type. Below
is an example minimal Appfile specifying the "my-app" application name and "go"
project type:

    application {
	name = "my-app"
	type = "go"
    }

If you believe Otto should've been able to automatically detect your
project type, then please open an issue with the Otto project.
`
