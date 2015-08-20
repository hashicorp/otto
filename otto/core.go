package otto

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
	"github.com/hashicorp/terraform/dag"
)

// Core is the main struct to use to interact with Otto as a library.
type Core struct {
	appfile         *appfile.File
	appfileCompiled *appfile.Compiled
	apps            map[app.Tuple]app.Factory
	dir             directory.Backend
	infras          map[string]infrastructure.Factory
	outputDir       string
	ui              ui.Ui
}

// CoreConfig is configuration for creating a new core with NewCore.
type CoreConfig struct {
	// OutputDir is the directory where data will be written. Each
	// compilation will clear this directory prior to writing to it.
	OutputDir string

	// Appfile is the appfile that this core will be using for configuration.
	// This must be a compiled Appfile.
	Appfile *appfile.Compiled

	// Directory is the directory where data is stored about this Appfile.
	Directory directory.Backend

	// Apps is the map of available app implementations.
	Apps map[app.Tuple]app.Factory

	// Infrastructures is the map of available infrastructures. The
	// value is a factory that can create the infrastructure impl.
	Infrastructures map[string]infrastructure.Factory

	// Ui is the Ui that will be used to communicate with the user.
	Ui ui.Ui
}

// NewCore creates a new core.
//
// Once this function is called, this CoreConfig should not be used again
// or modified, since the Core may use parts of it without deep copying.
func NewCore(c *CoreConfig) (*Core, error) {
	return &Core{
		appfile:         c.Appfile.File,
		appfileCompiled: c.Appfile,
		apps:            c.Apps,
		dir:             c.Directory,
		infras:          c.Infrastructures,
		outputDir:       c.OutputDir,
		ui:              c.Ui,
	}, nil
}

// Compile takes the Appfile and compiles all the resulting data.
func (c *Core) Compile() error {
	// Get the infra implementation for this
	infra, infraCtx, err := c.infra()
	if err != nil {
		return err
	}

	// Delete the prior output directory
	log.Printf("[INFO] deleting prior compilation contents: %s", c.outputDir)
	if err := os.RemoveAll(c.outputDir); err != nil {
		return err
	}

	// Compile the infrastructure for our application
	log.Printf("[INFO] running infra compile...")
	if _, err := infra.Compile(infraCtx); err != nil {
		return err
	}

	// Walk through the dependencies and compile all of them.
	// We have to compile every dependency for dev building.
	err = c.walk(func(app app.App, ctx *app.Context, root bool) error {
		if !root {
			c.ui.Message(fmt.Sprintf(
				"Compiling dependency '%s'...",
				ctx.Appfile.Application.Name))
		} else {
			c.ui.Message(fmt.Sprintf(
				"Compiling main application..."))
		}

		_, err := app.Compile(ctx)
		return err
	})

	return nil
}

func (c *Core) walk(f func(app.App, *app.Context, bool) error) error {
	root, err := c.appfileCompiled.Graph.Root()
	if err != nil {
		return fmt.Errorf(
			"Error loading app: %s", err)
	}

	// Walk the appfile graph.
	var stop int32 = 0
	return c.appfileCompiled.Graph.Walk(func(raw dag.Vertex) (err error) {
		// If we're told to stop (something else had an error), then stop early.
		// Graphs walks by default will complete all disjoint parts of the
		// graph before failing, but Otto doesn't have to do that.
		if atomic.LoadInt32(&stop) != 0 {
			return nil
		}

		// If we exit with an error, then mark the stop atomic
		defer func() {
			if err != nil {
				atomic.StoreInt32(&stop, 1)
			}
		}()

		// Convert to the rich vertex type so that we can access data
		v := raw.(*appfile.CompiledGraphVertex)

		// Get the context and app for this appfile
		appCtx, err := c.appContext(v.File)
		if err != nil {
			return fmt.Errorf(
				"Error loading Appfile for '%s': %s",
				dag.VertexName(raw), err)
		}
		app, err := c.app(appCtx)
		if err != nil {
			return fmt.Errorf(
				"Error loading App implementation for '%s': %s",
				dag.VertexName(raw), err)
		}

		// Call our callback
		return f(app, appCtx, raw == root)
	})
}

// Dev starts a dev environment for the current application. For destroying
// and other tasks against the dev environment, use the generic `Execute`
// method.
func (c *Core) Dev() error {
	// First get the root app and context, since we need that for
	// everything else.
	root, err := c.appfileCompiled.Graph.Root()
	if err != nil {
		return fmt.Errorf(
			"Error loading app: %s", err)
	}
	appCtx, err := c.appContext(root.(*appfile.CompiledGraphVertex).File)
	if err != nil {
		return fmt.Errorf("Error loading Appfile: %s", err)
	}
	app, err := c.app(appCtx)
	if err != nil {
		return fmt.Errorf("Error loading App: %s", err)
	}

	// Go through all the dependencies and build their immutable
	// dev environment pieces for the final configuration.
	var stop int32 = 0
	err = c.appfileCompiled.Graph.Walk(func(raw dag.Vertex) (err error) {
		// If this is the root, skip it. We do that separately.
		if raw == root {
			return nil
		}

		// If we're told to stop (something else had an error), then stop early
		if atomic.LoadInt32(&stop) != 0 {
			return nil
		}

		// If we exit with an error, then mark the stop atomic
		defer func() {
			if err != nil {
				atomic.StoreInt32(&stop, 1)
			}
		}()

		// Convert to the rich vertex type so that we can access data
		v := raw.(*appfile.CompiledGraphVertex)

		// Get the context and app for this appfile
		depCtx, err := c.appContext(v.File)
		if err != nil {
			return fmt.Errorf(
				"Error loading Appfile for dependency '%s': %s",
				dag.VertexName(raw), err)
		}
		dep, err := c.app(depCtx)
		if err != nil {
			return fmt.Errorf(
				"Error loading App implementation for dependency '%s': %s",
				dag.VertexName(raw), err)
		}

		// Build the development dependency
		println(dep)

		return nil
	})
	if err != nil {
		return err
	}

	println(app)
	return nil
}

// Execute executes the given task for this Appfile.
func (c *Core) Execute(opts *ExecuteOpts) error {
	switch opts.Task {
	case ExecuteTaskDev:
		return c.executeApp(opts)
	case ExecuteTaskInfra:
		return c.executeInfra(opts)
	default:
		return fmt.Errorf("unknown task: %s", opts.Task)
	}
}

func (c *Core) executeApp(opts *ExecuteOpts) error {
	// Get the infra implementation for this
	appCtx, err := c.appContext(c.appfile)
	if err != nil {
		return err
	}
	app, err := c.app(appCtx)
	if err != nil {
		return err
	}

	// Set the action and action args
	appCtx.Action = opts.Action
	appCtx.ActionArgs = opts.Args

	// Build the infrastructure compilation context
	switch opts.Task {
	case ExecuteTaskDev:
		return app.Dev(appCtx)
	default:
		panic(fmt.Sprintf("uknown task: %s", opts.Task))
	}
}

func (c *Core) executeInfra(opts *ExecuteOpts) error {
	// Get the infra implementation for this
	infra, infraCtx, err := c.infra()
	if err != nil {
		return err
	}

	// Set the action and action args
	infraCtx.Action = opts.Action
	infraCtx.ActionArgs = opts.Args

	// Build the infrastructure compilation context
	return infra.Execute(infraCtx)
}

func (c *Core) appContext(f *appfile.File) (*app.Context, error) {
	// We need the configuration for the active infrastructure
	// so that we can build the tuple below
	config := f.ActiveInfrastructure()
	if config == nil {
		return nil, fmt.Errorf(
			"infrastructure not found in appfile: %s",
			f.Project.Infrastructure)
	}

	// The tuple we're looking for is the application type, the
	// infrastructure type, and the infrastructure flavor. Build that
	// tuple.
	tuple := app.Tuple{
		App:         f.Application.Type,
		Infra:       f.Project.Infrastructure,
		InfraFlavor: config.Flavor,
	}

	// The output directory for data. This is either the main app so
	// it goes directly into "app" or it is a dependency and goes into
	// a dep folder.
	outputDir := filepath.Join(c.outputDir, "app")
	if id := f.ID(); id != c.appfile.ID() {
		outputDir = filepath.Join(
			c.outputDir, fmt.Sprintf("dep-%s", id))
	}

	return &app.Context{
		Dir:         outputDir,
		Tuple:       tuple,
		Appfile:     f,
		Application: f.Application,
		Ui:          c.ui,
	}, nil
}

func (c *Core) app(ctx *app.Context) (app.App, error) {
	// Look for the app impl. factory
	f, ok := c.apps[ctx.Tuple]
	if !ok {
		return nil, fmt.Errorf(
			"app implementation for tuple not found: %s", ctx.Tuple)
	}

	// Start the impl.
	result, err := f()
	if err != nil {
		return nil, fmt.Errorf(
			"app failed to start properly: %s", err)
	}

	return result, nil
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
		Dir:       outputDir,
		Infra:     config,
		Ui:        c.ui,
		Directory: c.dir,
	}, nil
}
