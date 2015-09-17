package terraform

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/foundation"
)

// Foundation is a helper for various operations a foundation must
// perform with Terraform.
type Foundation struct {
	// Dir is the directory where Terraform is run. If this isn't set, it'll
	// default to "#{ctx.Dir}/deploy".
	Dir string
}

// Infra manages a foundation using Terraform.
//
// This will verify the infrastruction is created and use that information
// to execute Terraform with the options given.
//
// This function implements foundation.Foundation.Infra.
func (f *Foundation) Infra(ctx *foundation.Context) error {
	switch ctx.Action {
	case "":
		if err := f.execute(ctx, "get", "."); err != nil {
			return err
		}

		return f.execute(ctx, "apply")
	case "destroy":
		return f.execute(ctx, "destroy", "-force")
	default:
		return fmt.Errorf("unknown action: %s", ctx.Action)
	}
}

func (f *Foundation) execute(ctx *foundation.Context, args ...string) error {
	project, err := Project(&ctx.Shared)
	if err != nil {
		return err
	}

	appInfra := ctx.Appfile.ActiveInfrastructure()

	// Foundations themselves are represented as infrastructure in the
	// backend. Let's look that up. If it doesn't exist, we have to create
	// it in order to get our UUID for storing state.
	lookup := directory.Lookup{Infra: appInfra.Name, Foundation: ctx.Tuple.Type}
	foundationInfra, err := ctx.Directory.GetInfra(&directory.Infra{Lookup: lookup})
	if err != nil {
		return fmt.Errorf(
			"Error looking up existing infrastructure data: %s\n\n"+
				"These errors are usually transient and can be fixed by retrying\n"+
				"the command. Additional causes of errors are networking or disk\n"+
				"issues that can be resolved external to Otto.",
			err)
	}
	if foundationInfra == nil {
		// If we don't have an infra, create one
		foundationInfra = &directory.Infra{Lookup: lookup}
		foundationInfra.State = directory.InfraStatePartial

		// Put the infrastructure so we can get the UUID to use for our state
		if err := ctx.Directory.PutInfra(foundationInfra); err != nil {
			return fmt.Errorf(
				"Error preparing infrastructure: %s\n\n"+
					"These errors are usually transient and can be fixed by retrying\n"+
					"the command. Additional causes of errors are networking or disk\n"+
					"issues that can be resolved external to Otto.",
				err)
		}
	}

	// Get the infrastructure state. The infrastructure must be
	// created for us to deploy to it.
	infra, err := ctx.Directory.GetInfra(&directory.Infra{
		Lookup: directory.Lookup{Infra: appInfra.Name}})
	if err != nil {
		return err
	}
	if infra == nil || infra.State != directory.InfraStateReady {
		return fmt.Errorf(
			"Infrastructure for this application hasn't been built yet.\n" +
				"Building a foundation requires a target infrastruction to\n" +
				"be built. Please run `otta infra` to build the underlying\n" +
				"infrastructure.")
	}

	// Construct the variables for Terraform from our queried infra
	vars := make(map[string]string)
	for k, v := range infra.Outputs {
		vars[k] = v
	}
	for k, v := range ctx.InfraCreds {
		vars[k] = v
	}

	// Get the directory
	tfDir := f.Dir
	if tfDir == "" {
		tfDir = filepath.Join(ctx.Dir, "deploy")
	}

	// Run Terraform!
	tf := &Terraform{
		Path:      project.Path(),
		Dir:       tfDir,
		Ui:        ctx.Ui,
		Variables: vars,
		Directory: ctx.Directory,
		StateId:   foundationInfra.ID,
	}
	if err := tf.Execute(args...); err != nil {
		return fmt.Errorf(
			"Error running Terraform: %s\n\n"+
				"Terraform usually has helpful error messages. Please read the error\n"+
				"messages above and resolve them. Sometimes simply re-running the\n"+
				"command again will work.",
			err)
	}

	return nil
}
