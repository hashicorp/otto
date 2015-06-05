package otto

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
)

// Core is the main struct to use to interact with Otto as a library.
type Core struct {
	appfile   *appfile.File
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
