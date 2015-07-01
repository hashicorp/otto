package otto

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
)

// Core is the main struct to use to interact with Otto as a library.
type Core struct {
	appfile   *appfile.File
	apps      map[app.Tuple]app.Factory
	infras    map[string]infrastructure.Factory
	outputDir string
	ui        ui.Ui
}

// CoreConfig is configuration for creating a new core with NewCore.
type CoreConfig struct {
	// OutputDir is the directory where data will be written. Each
	// compilation will clear this directory prior to writing to it.
	OutputDir string

	// Appfile is the appfile that this core will be using for configuration.
	Appfile *appfile.File

	// Apps is the map of available app implementations.
	Apps map[app.Tuple]app.Factory

	// Infrastructures is the map of available infrastructures. The
	// value is a factory that can create the infrastructure impl.
	Infrastructures map[string]infrastructure.Factory

	// Ui is the Ui that will be used to comunicate with the user.
	Ui ui.Ui
}

// NewCore creates a new core.
//
// Once this function is called, this CoreConfig should not be used again
// or modified, since the Core may use parts of it without deep copying.
func NewCore(c *CoreConfig) (*Core, error) {
	return &Core{
		appfile:   c.Appfile,
		apps:      c.Apps,
		infras:    c.Infrastructures,
		outputDir: c.OutputDir,
		ui:        c.Ui,
	}, nil
}

// Compile takes the Appfile and compiles all the resulting data.
func (c *Core) Compile() error {
	// Get the infra implementation for this
	infra, infraCtx, err := c.infra()
	if err != nil {
		return err
	}

	// Build the infrastructure compilation context
	_, err = infra.Compile(infraCtx)
	return err
}

// Execute executes the given task for this Appfile.
func (c *Core) Execute(opts *ExecuteOpts) error {
	switch opts.Task {
	case ExecuteTaskInfra:
		return c.executeInfra(opts)
	default:
		return fmt.Errorf("unknown task: %s", opts.Task)
	}
}

func (c *Core) executeInfra(opts *ExecuteOpts) error {
	// Get the infra implementation for this
	infra, infraCtx, err := c.infra()
	if err != nil {
		return err
	}

	// Build the infrastructure compilation context
	return infra.Execute(infraCtx)
}

func (c *Core) app() (app.App, *app.Context, error) {
	// We need the configuration for the active infrastructure
	// so that we can build the tuple below
	config := c.appfile.ActiveInfrastructure()
	if config == nil {
		return nil, nil, fmt.Errorf(
			"infrastructure not found in appfile: %s",
			c.appfile.Project.Infrastructure)
	}

	// The tuple we're looking for is the application type, the
	// infrastructure type, and the infrastructure flavor. Build that
	// tuple.
	tuple := app.Tuple{
		App:         c.appfile.Application.Type,
		Infra:       c.appfile.Project.Infrastructure,
		InfraFlavor: config.Flavor,
	}

	// Look for the app impl. factory
	f, ok := c.apps[tuple]
	if !ok {
		return nil, nil, fmt.Errorf(
			"app implementation for tuple not found: %s", tuple)
	}

	// Start the impl.
	result, err := f()
	if err != nil {
		return nil, nil, fmt.Errorf(
			"app failed to start properly: %s", err)
	}

	// The output directory for data
	outputDir := filepath.Join(c.outputDir, "app")

	return result, &app.Context{
		Dir:         outputDir,
		Application: c.appfile.Application,
	}, nil
}

func (c *Core) infra() (infrastructure.Infrastructure, *infrastructure.Context, error) {
	// Get the infrastructure factory
	f, ok := c.infras[c.appfile.Project.Infrastructure]
	if !ok {
		return nil, nil, fmt.Errorf(
			"infrastructure type not supported: %s",
			c.appfile.Project.Infrastructure)
	}

	// Get the infrastructure configuration
	config := c.appfile.ActiveInfrastructure()
	if config == nil {
		return nil, nil, fmt.Errorf(
			"infrastructure not found in appfile: %s",
			c.appfile.Project.Infrastructure)
	}

	// Start the infrastructure implementation
	infra, err := f()
	if err != nil {
		return nil, nil, err
	}

	// The output directory for data
	outputDir := filepath.Join(
		c.outputDir, fmt.Sprintf("infra-%s", c.appfile.Project.Infrastructure))

	// Build the context
	return infra, &infrastructure.Context{
		Dir:   outputDir,
		Infra: config,
		Ui:    c.ui,
	}, nil
}
