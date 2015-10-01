package otto

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/helper/localaddr"
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
	foundationMap   map[foundation.Tuple]foundation.Factory
	dataDir         string
	localDir        string
	compileDir      string
	ui              ui.Ui
}

// CoreConfig is configuration for creating a new core with NewCore.
type CoreConfig struct {
	// DataDir is the directory where local data will be stored that
	// is global to all Otto processes.
	//
	// LocalDir is the directory where data local to this single Appfile
	// will be stored. This isn't necessarilly cleared for compilation.
	//
	// CompiledDir is the directory where compiled data will be written.
	// Each compilation will clear this directory.
	DataDir    string
	LocalDir   string
	CompileDir string

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

	// Foundations is the map of available foundations. The
	// value is a factory that can create the impl.
	Foundations map[foundation.Tuple]foundation.Factory

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
		foundationMap:   c.Foundations,
		dataDir:         c.DataDir,
		localDir:        c.LocalDir,
		compileDir:      c.CompileDir,
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

	// Get all the foundation implementations (which are tied as singletons
	// to the infrastructure).
	foundations, foundationCtxs, err := c.foundations()
	if err != nil {
		return err
	}

	// Delete the prior output directory
	log.Printf("[INFO] deleting prior compilation contents: %s", c.compileDir)
	if err := os.RemoveAll(c.compileDir); err != nil {
		return err
	}

	// Compile the infrastructure for our application
	log.Printf("[INFO] running infra compile...")
	c.ui.Message("Compiling infra...")
	if _, err := infra.Compile(infraCtx); err != nil {
		return err
	}

	// Compile the foundation (not tied to any app). This compilation
	// of the foundation is used for `otto infra` to set everything up.
	log.Printf("[INFO] running foundation compilations")
	for i, f := range foundations {
		ctx := foundationCtxs[i]
		c.ui.Message(fmt.Sprintf(
			"Compiling foundation: %s", ctx.Tuple.Type))
		if _, err := f.Compile(ctx); err != nil {
			return err
		}
	}

	// Walk through the dependencies and compile all of them.
	// We have to compile every dependency for dev building.
	var resultLock sync.Mutex
	results := make([]*app.CompileResult, 0, len(c.appfileCompiled.Graph.Vertices()))
	err = c.walk(func(app app.App, ctx *app.Context, root bool) error {
		if !root {
			c.ui.Header(fmt.Sprintf(
				"Compiling dependency '%s'...",
				ctx.Appfile.Application.Name))
		} else {
			c.ui.Header(fmt.Sprintf(
				"Compiling main application..."))
		}

		// If this is the root, we set the dev dep fragments.
		if root {
			// We grab the lock just in case although if we're the
			// root this should be serialized.
			resultLock.Lock()
			ctx.DevDepFragments = make([]string, 0, len(results))
			for _, result := range results {
				if result.DevDepFragmentPath != "" {
					ctx.DevDepFragments = append(
						ctx.DevDepFragments, result.DevDepFragmentPath)
				}
			}
			resultLock.Unlock()
		}

		// Compile!
		result, err := app.Compile(ctx)
		if err != nil {
			return err
		}

		// Compile the foundations for this app
		subdirs := []string{"app-dev", "app-dev-dep", "app-build", "app-deploy"}
		for i, f := range foundations {
			fCtx := foundationCtxs[i]
			fCtx.Dir = ctx.FoundationDirs[i]
			if result != nil {
				fCtx.AppConfig = &result.FoundationConfig
			}

			if _, err := f.Compile(fCtx); err != nil {
				return err
			}

			// Make sure the subdirs exist
			for _, dir := range subdirs {
				if err := os.MkdirAll(filepath.Join(fCtx.Dir, dir), 0755); err != nil {
					return err
				}
			}
		}

		// Store the compilation result for later
		resultLock.Lock()
		defer resultLock.Unlock()
		results = append(results, result)

		return nil
	})

	return err
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

// Build builds the deployable artifact for the currently compiled
// Appfile.
func (c *Core) Build() error {
	// Get the infra implementation for this
	infra, infraCtx, err := c.infra()
	if err != nil {
		return err
	}
	if err := c.creds(infra, infraCtx); err != nil {
		return err
	}

	// We only use the root application for this task, upstream dependencies
	// don't have an effect on the build process.
	root, err := c.appfileCompiled.Graph.Root()
	if err != nil {
		return err
	}
	rootCtx, err := c.appContext(root.(*appfile.CompiledGraphVertex).File)
	if err != nil {
		return fmt.Errorf(
			"Error loading App: %s", err)
	}
	rootApp, err := c.app(rootCtx)
	if err != nil {
		return fmt.Errorf(
			"Error loading App: %s", err)
	}

	// Just update our shared data so we get the creds
	rootCtx.Shared.InfraCreds = infraCtx.Shared.InfraCreds

	return rootApp.Build(rootCtx)
}

// Deploy deploys the application.
//
// Deploy supports subactions, which can be specified with action and args.
// Action can be "" to get the default deploy behavior.
func (c *Core) Deploy(action string, args []string) error {
	// Get the infra implementation for this
	infra, infraCtx, err := c.infra()
	if err != nil {
		return err
	}

	// Special case: don't try to fetch creds during `help` or `info`
	if action != "help" && action != "info" {
		if err := c.creds(infra, infraCtx); err != nil {
			return err
		}
	}

	// TODO: Verify that upstream dependencies are deployed

	// We only use the root application for this task, upstream dependencies
	// don't have an effect on the build process.
	root, err := c.appfileCompiled.Graph.Root()
	if err != nil {
		return err
	}
	rootCtx, err := c.appContext(root.(*appfile.CompiledGraphVertex).File)
	if err != nil {
		return fmt.Errorf(
			"Error loading App: %s", err)
	}
	rootApp, err := c.app(rootCtx)
	if err != nil {
		return fmt.Errorf(
			"Error loading App: %s", err)
	}

	// Update our shared data so we get the creds
	rootCtx.Shared.InfraCreds = infraCtx.Shared.InfraCreds

	// Pass through the requested action
	rootCtx.Action = action
	rootCtx.ActionArgs = args

	return rootApp.Deploy(rootCtx)
}

// Dev starts a dev environment for the current application. For destroying
// and other tasks against the dev environment, use the generic `Execute`
// method.
func (c *Core) Dev() error {
	// We need to get the root data separately since we need that for
	// all the function calls into the dependencies.
	root, err := c.appfileCompiled.Graph.Root()
	if err != nil {
		return err
	}
	rootCtx, err := c.appContext(root.(*appfile.CompiledGraphVertex).File)
	if err != nil {
		return fmt.Errorf(
			"Error loading App: %s", err)
	}
	rootApp, err := c.app(rootCtx)
	if err != nil {
		return fmt.Errorf(
			"Error loading App: %s", err)
	}

	// Go through all the dependencies and build their immutable
	// dev environment pieces for the final configuration.
	err = c.walk(func(appImpl app.App, ctx *app.Context, root bool) error {
		// If it is the root, we just return and do nothing else since
		// the root is a special case where we're building the actual
		// dev environment.
		if root {
			return nil
		}

		// Get the path to where we'd cache the dependency if we have
		// cached it...
		cachePath := filepath.Join(ctx.CacheDir, "dev-dep.json")

		// Check if we've cached this. If so, then use the cache.
		if _, err := app.ReadDevDep(cachePath); err == nil {
			ctx.Ui.Header(fmt.Sprintf(
				"Using cached dev dependency for '%s'",
				ctx.Appfile.Application.Name))
			return nil
		}

		// Build the development dependency
		dep, err := appImpl.DevDep(rootCtx, ctx)
		if err != nil {
			return fmt.Errorf(
				"Error building dependency for dev '%s': %s",
				ctx.Appfile.Application.Name,
				err)
		}

		// If we have a dependency with files, then verify the files
		// and store it in our cache directory so we can retrieve it
		// later.
		if dep != nil && len(dep.Files) > 0 {
			if err := dep.RelFiles(ctx.CacheDir); err != nil {
				return fmt.Errorf(
					"Error caching dependency for dev '%s': %s",
					ctx.Appfile.Application.Name,
					err)
			}

			if err := app.WriteDevDep(cachePath, dep); err != nil {
				return fmt.Errorf(
					"Error caching dependency for dev '%s': %s",
					ctx.Appfile.Application.Name,
					err)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// All the development dependencies are built/loaded. We now have
	// everything we need to build the complete development environment.
	return rootApp.Dev(rootCtx)
}

// Infra manages the infrastructure for this Appfile.
//
// Infra supports subactions, which can be specified with action and args.
// Infra recognizes two special actions: "" (blank string) and "destroy".
// The former expects to create or update the complete infrastructure,
// and the latter will destroy the infrastructure.
func (c *Core) Infra(action string, args []string) error {
	// Get the infra implementation for this
	infra, infraCtx, err := c.infra()
	if err != nil {
		return err
	}
	if action == "" || action == "destroy" {
		if err := c.creds(infra, infraCtx); err != nil {
			return err
		}
	}

	// Set the action and action args
	infraCtx.Action = action
	infraCtx.ActionArgs = args

	// If we need the foundations, then get them
	var foundations []foundation.Foundation
	var foundationCtxs []*foundation.Context
	if action == "" || action == "destroy" {
		foundations, foundationCtxs, err = c.foundations()
		if err != nil {
			return err
		}
	}

	// If we're doing anything other than destroying, then
	// run the execution now.
	if action != "destroy" {
		if err := infra.Execute(infraCtx); err != nil {
			return err
		}
	}

	// If we have any foundations, we now run their infra deployment.
	// This should only ever execute if action is to deploy or destroy,
	// since those are the only cases that we load foundations.
	for i, f := range foundations {
		ctx := foundationCtxs[i]
		ctx.Action = action
		ctx.ActionArgs = args
		ctx.InfraCreds = infraCtx.InfraCreds

		log.Printf(
			"[INFO] infra action '%s' on foundation '%s'",
			action, ctx.Tuple.Type)

		switch action {
		case "":
			c.ui.Header(fmt.Sprintf(
				"Building infrastructure for foundation: %s",
				ctx.Tuple.Type))
		case "destroy":
			c.ui.Header(fmt.Sprintf(
				"Destroying infrastructure for foundation: %s",
				ctx.Tuple.Type))
		}

		if err := f.Infra(ctx); err != nil {
			return err
		}
	}

	// If the action is destroy, we run the infrastructure execution
	// here. We mirror creation above since in the destruction case
	// we need to first destroy all applications and foundations that
	// are using this infra.
	if action == "destroy" {
		if err := infra.Execute(infraCtx); err != nil {
			return err
		}
	}

	// Output the right thing
	switch action {
	case "":
		infraCtx.Ui.Header("[green]Infrastructure successfully created!")
		infraCtx.Ui.Message(
			"[green]The infrastructure necessary to deploy this application\n" +
				"is now available. You can now deploy using `otto deploy`.")
	case "destroy":
		infraCtx.Ui.Header("[green]Infrastructure successfully destroyed!")
		infraCtx.Ui.Message(
			"[green]The infrastructure necessary to run this application and\n" +
				"all other applications in this project has been destroyed.")
	}

	return nil
}

// Status outputs to the UI the status of all the stages of this application.
func (c *Core) Status() error {
	// Start loading the status info in a goroutine
	statusCh := make(chan *statusInfo, 1)
	go c.statusInfo(statusCh)

	// Wait for the status. If this takes longer than a certain amount
	// of time then we show a loading message.
	var status *statusInfo
	select {
	case status = <-statusCh:
	case <-time.After(150 * time.Millisecond):
		c.ui.Header("Loading status...")
		c.ui.Message(fmt.Sprintf(
			"Depending on your configured directory backend, this may require\n" +
				"network operations and can take some time. On a typical broadband\n" +
				"connection, this shouldn't take more than a few seconds."))
	}
	if status == nil {
		status = <-statusCh
	}

	// Create the status texts
	devStatus := "[reset]NOT CREATED"
	if status.Dev.IsReady() {
		devStatus = "[green]CREATED"
	}
	buildStatus := "[reset]NOT BUILT"
	if status.Build != nil {
		buildStatus = "[green]BUILD READY"
	}
	deployStatus := "[reset]NOT DEPLOYED"
	if status.Deploy.IsDeployed() {
		deployStatus = "[green]DEPLOYED"
	} else if status.Deploy.IsFailed() {
		deployStatus = "[reset]DEPLOY FAILED"
	}
	infraStatus := "[reset]NOT CREATED"
	if status.Infra.IsReady() {
		infraStatus = "[green]READY"
	} else if status.Infra.IsPartial() {
		infraStatus = "[yellow]PARTIAL"
	}

	// Get the active infra
	infra := c.appfile.ActiveInfrastructure()

	c.ui.Header("App Info")
	c.ui.Message(fmt.Sprintf(
		"Application:    %s (%s)",
		c.appfile.Application.Name, c.appfile.Application.Type))
	c.ui.Message(fmt.Sprintf("Project:        %s", c.appfile.Project.Name))
	c.ui.Message(fmt.Sprintf(
		"Infrastructure: %s (%s)",
		infra.Type, infra.Flavor))

	c.ui.Header("Component Status")
	c.ui.Message(fmt.Sprintf("Dev environment: %s", devStatus))
	c.ui.Message(fmt.Sprintf("Infra:           %s", infraStatus))
	c.ui.Message(fmt.Sprintf("Build:           %s", buildStatus))
	c.ui.Message(fmt.Sprintf("Deploy:          %s", deployStatus))

	return nil
}

// Execute executes the given task for this Appfile.
func (c *Core) Execute(opts *ExecuteOpts) error {
	switch opts.Task {
	case ExecuteTaskDev:
		return c.executeApp(opts)
	default:
		return fmt.Errorf("unknown task: %s", opts.Task)
	}
}

// creds reads the credentials if we have them, or queries the user
// for infrastructure credentials using the infrastructure if we
// don't have them.
func (c *Core) creds(
	infra infrastructure.Infrastructure,
	infraCtx *infrastructure.Context) error {
	// Output to the user some information about what is about to
	// happen here...
	infraCtx.Ui.Header(fmt.Sprintf(
		"Detecting infrastructure credentials for: %s (%s)",
		infraCtx.Infra.Name, infraCtx.Infra.Type))

	// The path to where we put the encrypted creds
	path := filepath.Join(c.dataDir, "cache", "creds", infraCtx.Infra.Name)

	// Determine whether we believe the creds exist already or not
	var exists bool
	if _, err := os.Stat(path); err == nil {
		exists = true
	} else {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
	}

	var creds map[string]string
	if exists {
		infraCtx.Ui.Message(
			"Cached and encrypted infrastructure credentials found.\n" +
				"Otto will now ask you for the password to decrypt these\n" +
				"credentials.\n\n")

		// If they exist, ask for the password
		value, err := infraCtx.Ui.Input(&ui.InputOpts{
			Id:          "creds_password",
			Query:       "Encrypted Credentials Password",
			Description: strings.TrimSpace(credsQueryPassExists),
			Hide:        true,
			EnvVars:     []string{"OTTO_CREDS_PASSWORD"},
		})
		if err != nil {
			return err
		}

		// If the password is not blank, then just read the credentials
		if value != "" {
			plaintext, err := cryptRead(path, value)
			if err == nil {
				err = json.Unmarshal(plaintext, &creds)
			}
			if err != nil {
				return fmt.Errorf(
					"error reading encrypted credentials: %s\n\n"+
						"If this error persists, you can force Otto to ask for credentials\n"+
						"again by inputting the empty password as the password.",
					err)
			}
		}
	}

	// If we don't have creds, then we need to query the user via
	// the infrastructure implementation.
	if creds == nil {
		infraCtx.Ui.Message(
			"Existing infrastructure credentials were not found! Otto will\n" +
				"now ask you for infrastructure credentials. These will be encrypted\n" +
				"and saved on disk so this doesn't need to be repeated.\n\n" +
				"IMPORTANT: If you're re-entering new credentials, make sure the\n" +
				"credentials are for the same account, otherwise you may lose\n" +
				"access to your existing infrastructure Otto set up.\n\n")

		var err error
		creds, err = infra.Creds(infraCtx)
		if err != nil {
			return err
		}

		// Now that we have the credentials, we need to ask for the
		// password to encrypt and store them.
		var password string
		for password == "" {
			password, err = infraCtx.Ui.Input(&ui.InputOpts{
				Id:          "creds_password",
				Query:       "Password for Encrypting Credentials",
				Description: strings.TrimSpace(credsQueryPassNew),
				Hide:        true,
				EnvVars:     []string{"OTTO_CREDS_PASSWORD"},
			})
			if err != nil {
				return err
			}
		}

		// With the password, encrypt and write the data
		plaintext, err := json.Marshal(creds)
		if err != nil {
			// creds is a map[string]string, so this shouldn't ever fail
			panic(err)
		}

		if err := cryptWrite(path, password, plaintext); err != nil {
			return fmt.Errorf(
				"error writing encrypted credentials: %s", err)
		}
	}

	// Set the credentials
	infraCtx.InfraCreds = creds

	// Let the infrastructure do whatever it likes to verify that the credentials
	// are good, so we can fail fast in case there's a problem.
	if err := infra.VerifyCreds(infraCtx); err != nil {
		return err
	}

	return nil
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
		Infra:       config.Type,
		InfraFlavor: config.Flavor,
	}

	// The output directory for data. This is either the main app so
	// it goes directly into "app" or it is a dependency and goes into
	// a dep folder.
	outputDir := filepath.Join(c.compileDir, "app")
	if id := f.ID; id != c.appfile.ID {
		outputDir = filepath.Join(
			c.compileDir, fmt.Sprintf("dep-%s", id))
	}

	// The cache directory for this app
	cacheDir := filepath.Join(c.dataDir, "cache", f.ID)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, fmt.Errorf(
			"error making cache directory '%s': %s",
			cacheDir, err)
	}

	// Build the contexts for the foundations. We use this
	// to also compile the list of foundation dirs.
	foundationDirs := make([]string, len(config.Foundations))
	for i, f := range config.Foundations {
		foundationDirs[i] = filepath.Join(
			outputDir, fmt.Sprintf("foundation-%s", f.Name))
	}

	// Get the dev IP address
	ipDB := &localaddr.CachedDB{
		DB:        &localaddr.DB{Path: filepath.Join(c.dataDir, "ip.db")},
		CachePath: filepath.Join(c.localDir, "dev_ip"),
	}
	ip, err := ipDB.IP()
	if err != nil {
		return nil, fmt.Errorf(
			"Error retrieving dev IP address: %s", err)
	}

	return &app.Context{
		Dir:          outputDir,
		CacheDir:     cacheDir,
		LocalDir:     c.localDir,
		Tuple:        tuple,
		Application:  f.Application,
		DevIPAddress: ip.String(),
		Shared: context.Shared{
			Appfile:        f,
			FoundationDirs: foundationDirs,
			InstallDir:     filepath.Join(c.dataDir, "binaries"),
			Directory:      c.dir,
			Ui:             c.ui,
		},
	}, nil
}

func (c *Core) app(ctx *app.Context) (app.App, error) {
	log.Printf("[INFO] Loading app implementation for Tuple: %s", ctx.Tuple)

	// Look for the app impl. factory
	f := app.TupleMap(c.apps).Lookup(ctx.Tuple)
	if f == nil {
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
	// Get the infrastructure configuration
	config := c.appfile.ActiveInfrastructure()
	if config == nil {
		return nil, nil, fmt.Errorf(
			"infrastructure not found in appfile: %s",
			c.appfile.Project.Infrastructure)
	}

	// Get the infrastructure factory
	f, ok := c.infras[config.Type]
	if !ok {
		return nil, nil, fmt.Errorf(
			"infrastructure type not supported: %s",
			config.Type)
	}

	// Start the infrastructure implementation
	infra, err := f()
	if err != nil {
		return nil, nil, err
	}

	// The output directory for data
	outputDir := filepath.Join(
		c.compileDir, fmt.Sprintf("infra-%s", c.appfile.Project.Infrastructure))

	// Build the context
	return infra, &infrastructure.Context{
		Dir:   outputDir,
		Infra: config,
		Shared: context.Shared{
			Appfile:    c.appfile,
			InstallDir: filepath.Join(c.dataDir, "binaries"),
			Directory:  c.dir,
			Ui:         c.ui,
		},
	}, nil
}

func (c *Core) foundations() ([]foundation.Foundation, []*foundation.Context, error) {
	// Get the infrastructure configuration
	config := c.appfile.ActiveInfrastructure()
	if config == nil {
		return nil, nil, fmt.Errorf(
			"infrastructure not found in appfile: %s",
			c.appfile.Project.Infrastructure)
	}

	// If there are no foundations, return nothing.
	if len(config.Foundations) == 0 {
		return nil, nil, nil
	}

	// Create the arrays for our list
	fs := make([]foundation.Foundation, 0, len(config.Foundations))
	ctxs := make([]*foundation.Context, 0, cap(fs))
	for _, f := range config.Foundations {
		// The tuple we're looking for is the foundation type, the
		// infrastructure type, and the infrastructure flavor. Build that
		// tuple.
		tuple := foundation.Tuple{
			Type:        f.Name,
			Infra:       config.Type,
			InfraFlavor: config.Flavor,
		}

		// Look for the matching foundation
		fun := foundation.TupleMap(c.foundationMap).Lookup(tuple)
		if fun == nil {
			return nil, nil, fmt.Errorf(
				"foundation implementation for tuple not found: %s",
				tuple)
		}

		// Instantiate the implementation
		impl, err := fun()
		if err != nil {
			return nil, nil, err
		}

		// The output directory for data
		outputDir := filepath.Join(
			c.compileDir, fmt.Sprintf("foundation-%s", f.Name))

		// Build the context
		ctx := &foundation.Context{
			Config: f.Config,
			Dir:    outputDir,
			Tuple:  tuple,
			Shared: context.Shared{
				Appfile:    c.appfile,
				InstallDir: filepath.Join(c.dataDir, "binaries"),
				Directory:  c.dir,
				Ui:         c.ui,
			},
		}

		// Add to our results
		fs = append(fs, impl)
		ctxs = append(ctxs, ctx)
	}

	return fs, ctxs, nil
}

const credsQueryPassExists = `
Infrastructure credentials are required for this operation. Otto found
saved credentials that are password protected. Please enter the password
to decrypt these credentials. You may also just hit <enter> and leave
the password blank to force Otto to ask for the credentials again.
`

const credsQueryPassNew = `
This password will be used to encrypt and save the credentials so they
don't need to be repeated multiple times.
`
