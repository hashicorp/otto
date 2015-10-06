package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/helper/bindata"
)

// AppOptions are the options for compiling an application.
//
// These options may be modified during customization processing, and
// in fact that is an intended use case and common pattern. To do this,
// use the AppCustomizationFunc method. See some of the builtin types for
// examples.
type AppOptions struct {
	// Ctx is the app context of this compilation.
	Ctx *app.Context

	// FoundationConfig is the configuration for the foundation that
	// will be returned as the compilation result.
	FoundationConfig foundation.Config

	// Bindata is the data that is used for templating. This must be set.
	// Template data should also be set on this. This will be modified with
	// default template data if those keys are not set.
	Bindata *bindata.Data

	// Customizations is a list of helpers to process customizations
	// in the Appfile. See the Customization docs for more information.
	Customizations []*Customization

	// Callbacks are called just prior to compilation completing.
	Callbacks []CompileCallback
}

// CompileCallback is a callback that can be registered to be run after
// compilation. To access any data within this callback, it should be created
// as a closure around the AppOptions.
type CompileCallback func() error

// App is an opinionated compilation function to help implement
// app.App.Compile.
//
// AppOptions may be modified by this function during this call.
func App(opts *AppOptions) (*app.CompileResult, error) {
	ctx := opts.Ctx

	// Setup the basic templating data. We put this into the "data" local
	// var just so that it is easier to reference.
	//
	// The exact default data put into the context is documented above.
	data := opts.Bindata
	if data.Context == nil {
		data.Context = make(map[string]interface{})
		opts.Bindata = data
	}

	data.Context["app_type"] = ctx.Appfile.Application.Type
	data.Context["name"] = ctx.Appfile.Application.Name
	data.Context["dev_fragments"] = ctx.DevDepFragments
	data.Context["dev_ip_address"] = ctx.DevIPAddress

	if data.Context["path"] == nil {
		data.Context["path"] = make(map[string]string)
	}
	pathMap := data.Context["path"].(map[string]string)
	pathMap["cache"] = ctx.CacheDir
	pathMap["compiled"] = ctx.Dir
	pathMap["working"] = filepath.Dir(ctx.Appfile.Path)
	foundationDirsContext := map[string][]string{
		"dev":     make([]string, len(ctx.FoundationDirs)),
		"dev_dep": make([]string, len(ctx.FoundationDirs)),
		"build":   make([]string, len(ctx.FoundationDirs)),
		"deploy":  make([]string, len(ctx.FoundationDirs)),
	}
	for i, dir := range ctx.FoundationDirs {
		foundationDirsContext["dev"][i] = filepath.Join(dir, "app-dev")
		foundationDirsContext["dev_dep"][i] = filepath.Join(dir, "app-dev-dep")
		foundationDirsContext["build"][i] = filepath.Join(dir, "app-build")
		foundationDirsContext["deploy"][i] = filepath.Join(dir, "app-deploy")
	}
	data.Context["foundation_dirs"] = foundationDirsContext

	// Setup the shared data
	if data.SharedExtends == nil {
		data.SharedExtends = make(map[string]*bindata.Data)
	}
	data.SharedExtends["compile"] = &bindata.Data{
		Asset:    Asset,
		AssetDir: AssetDir,
	}

	// Process the customizations!
	err := processCustomizations(&processOpts{
		Customizations: opts.Customizations,
		Appfile:        ctx.Appfile,
		Bindata:        data,
	})
	if err != nil {
		return nil, err
	}

	// Create the directory list that we'll copy from, and copy those
	// directly into the compilation directory.
	bindirs := []string{
		"data/common",
		fmt.Sprintf("data/%s-%s", ctx.Tuple.Infra, ctx.Tuple.InfraFlavor),
	}
	for _, dir := range bindirs {
		// Copy all the common files that exist
		if err := data.CopyDir(ctx.Dir, dir); err != nil {
			// Ignore any directories that don't exist
			if strings.Contains(err.Error(), "not found") {
				continue
			}

			return nil, err
		}
	}

	// Callbacks
	for _, cb := range opts.Callbacks {
		if err := cb(); err != nil {
			return nil, err
		}
	}

	// If the DevDep fragment exists, then use it
	fragmentPath := filepath.Join(ctx.Dir, "dev-dep", "Vagrantfile.fragment")
	if _, err := os.Stat(fragmentPath); err != nil {
		fragmentPath = ""
	}

	// Set some defaults here
	if opts.FoundationConfig.ServiceName == "" {
		opts.FoundationConfig.ServiceName = opts.Ctx.Application.Name
	}

	return &app.CompileResult{
		FoundationConfig:   opts.FoundationConfig,
		DevDepFragmentPath: fragmentPath,
	}, nil
}
