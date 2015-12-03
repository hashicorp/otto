package compile

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/helper/bindata"
)

// FoundationOptions are the options for compiling a foundation.
//
// These options may be modified during customization processing, and
// in fact that is an intended use case and common pattern. To do this,
// use the FoundationCustomizationFunc method. See some of the builtin types for
// examples.
type FoundationOptions struct {
	// Ctx is the foundation context of this compilation.
	Ctx *foundation.Context

	// Bindata is the data that is used for templating. This must be set.
	// Template data should also be set on this. This will be modified with
	// default template data if those keys are not set.
	Bindata *bindata.Data

	// Customization is used to configure the customizations for this
	// application. See the Customization type docs for more info.
	Customization *Customization

	// Callbacks are called just prior to compilation completing.
	Callbacks []CompileCallback
}

// Foundation is an opinionated compilation function to help implement
// foundation.Foundation.Compile.
//
// FoundationOptions may be modified by this function during this call.
func Foundation(opts *FoundationOptions) (*foundation.CompileResult, error) {
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
	data.Context["name"] = ctx.Appfile.Application.Name
	data.Context["path"] = map[string]string{
		"compiled": ctx.Dir,
		"working":  filepath.Dir(ctx.Appfile.Path),
	}
	data.Context["app_config"] = ctx.AppConfig
	if ctx.AppConfig == nil {
		data.Context["app_config"] = &foundation.Config{}
	}

	// Process the customizations!
	err := processCustomizations(
		ctx.Appfile.Customization,
		opts.Customization)
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

	return nil, nil
}
