package compile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/helper/bindata"
)

type AppOptions struct {
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
// compilation.
type CompileCallback func(*app.Context, *bindata.Data) error

// App is an opinionated compilation function to help implement
// app.App.Compile.
//
// AppOptions may be modified during this call.
func App(ctx *app.Context, opts *AppOptions) (*app.CompileResult, error) {
	// Setup the basic templating data. We put this into the "data" local
	// var just so that it is easier to reference.
	//
	// The exact default data put into the context is documented above.
	data := opts.Bindata
	if data.Context == nil {
		data.Context = make(map[string]interface{})
		opts.Bindata = data
	}
	data.Context["name"] = ctx.Appfile.Application.Name
	data.Context["dev_fragments"] = ctx.DevDepFragments
	data.Context["path"] = map[string]string{
		"cache":    ctx.CacheDir,
		"compiled": ctx.Dir,
		"working":  filepath.Dir(ctx.Appfile.Path),
	}

	// Process the customizations!
	customResult, err := processCustomizations(&processOpts{
		Customizations: opts.Customizations,
		Appfile:        ctx.Appfile,
		Bindata:        data,
	})
	if err != nil {
		return nil, err
	}

	// Store the callbacks
	callbacks := make([]CompileCallback, 0, len(opts.Callbacks)+len(customResult.Callbacks))
	callbacks = append(callbacks, opts.Callbacks...)
	callbacks = append(callbacks, customResult.Callbacks...)

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
	for _, cb := range callbacks {
		if err := cb(ctx, opts.Bindata); err != nil {
			return nil, err
		}
	}

	// If the DevDep fragment exists, then use it
	fragmentPath := filepath.Join(ctx.Dir, "dev-dep", "Vagrantfile.fragment")
	if _, err := os.Stat(fragmentPath); err != nil {
		fragmentPath = ""
	}

	return &app.CompileResult{
		DevDepFragmentPath: fragmentPath,
	}, nil
}
