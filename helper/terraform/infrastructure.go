package terraform

import (
	"fmt"

	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/infrastructure"
)

// Infrastructure implements infrastructure.Infrastructure and is a
// higher level framework for writing infrastructure implementations that
// use Terraform.
//
// This implementation will automatically:
//
//   * Save/restore state files via the directory service
//   * Populate infrastructure data in the directory (w/ Terraform outputs)
//   * Handle many edge case scenarios gracefully
//
type Infrastructure struct {
	// Creds is a function that gathers credentials. See helper/creds
	// for nice helpers for implementing this function.
	CredsFunc func(*infrastructure.Context) (map[string]string, error)

	// Bindata is the bindata.Data structure where assets can be found
	// for compilation. The data for the various flavors is expected to
	// live in "data/#{flavor}"
	Bindata *bindata.Data

	// Variables are additional variables to pass into Terraform.
	Variables map[string]string
}

func (i *Infrastructure) Creds(ctx *infrastructure.Context) (map[string]string, error) {
	return i.CredsFunc(ctx)
}

func (i *Infrastructure) Execute(ctx *infrastructure.Context) error {
	switch ctx.Action {
	case "destroy":
		return i.execute(ctx, "destroy")
	case "":
		return i.execute(ctx, "apply")
	default:
		return nil
	}
}

func (i *Infrastructure) execute(ctx *infrastructure.Context, command string) error {
	project, err := Project(&ctx.Shared)
	if err != nil {
		return err
	}

	// Build the variables
	vars := make(map[string]string)
	for k, v := range ctx.InfraCreds {
		vars[k] = v
	}
	for k, v := range i.Variables {
		vars[k] = v
	}

	// Setup the lookup information and query the existing infra so we
	// can get our UUID for storing data.
	lookup := directory.Lookup{Infra: ctx.Infra.Name}
	infra, err := ctx.Directory.GetInfra(&directory.Infra{Lookup: lookup})
	if err != nil {
		return fmt.Errorf(
			"Error looking up existing infrastructure data: %s\n\n"+
				"These errors are usually transient and can be fixed by retrying\n"+
				"the command. Additional causes of errors are networking or disk\n"+
				"issues that can be resolved external to Otto.",
			err)
	}
	if infra == nil {
		// If we don't have an infra, create one
		infra = &directory.Infra{Lookup: lookup}
		infra.State = directory.InfraStatePartial

		// Put the infrastructure so we can get the UUID to use for our state
		if err := ctx.Directory.PutInfra(infra); err != nil {
			return fmt.Errorf(
				"Error preparing infrastructure: %s\n\n"+
					"These errors are usually transient and can be fixed by retrying\n"+
					"the command. Additional causes of errors are networking or disk\n"+
					"issues that can be resolved external to Otto.",
				err)
		}
	}

	// Build our executor
	tf := &Terraform{
		Path:      project.Path(),
		Dir:       ctx.Dir,
		Ui:        ctx.Ui,
		Variables: vars,
		Directory: ctx.Directory,
		StateId:   infra.ID,
	}

	ctx.Ui.Header("Executing Terraform to manage infrastructure...")
	ctx.Ui.Message(
		"Raw Terraform output will begin streaming in below. Otto\n" +
			"does not create this output. It is mirrored directly from\n" +
			"Terraform while the infrastructure is being created.\n\n" +
			"Terraform may ask for input. For infrastructure provider\n" +
			"credentials, be sure to enter the same credentials\n" +
			"consistently within the same Otto environment." +
			"\n\n")

	// Start the Terraform command
	err = tf.Execute(command)
	if err != nil {
		err = fmt.Errorf("Error running Terraform: %s", err)
		infra.State = directory.InfraStatePartial
	}

	ctx.Ui.Header("Terraform execution complete. Saving results...")

	// Read the outputs if everything is looking good so far
	if err == nil {
		infra.State = directory.InfraStateReady
		infra.Outputs, err = tf.Outputs()
		if err != nil {
			err = fmt.Errorf("Error reading Terraform outputs: %s", err)
			infra.State = directory.InfraStatePartial
		}
	}

	// Save the infrastructure information
	if err := ctx.Directory.PutInfra(infra); err != nil {
		return fmt.Errorf(
			"Error storing infrastructure data: %s\n\n"+
				"This means that Otto won't be able to know that your infrastructure\n"+
				"was successfully created. Otto tries a few times to save the\n"+
				"infrastructure. At this point in time, Otto doesn't support gracefully\n"+
				"recovering from this error. Your infrastructure is now orphaned from\n"+
				"Otto's management. Please reference the community for help.\n\n"+
				"A future version of Otto will resolve this.",
			err)
	}

	// If there was an error during the process, then return that.
	if err != nil {
		return fmt.Errorf("Error reading Terraform outputs: %s\n\n"+
			"In this case, Otto is unable to consider the infrastructure ready.\n"+
			"Otto won't lose your infrastructure information. You may just need\n"+
			"to run `otto infra` again and it may work. If this problem persists,\n"+
			"please see the error message and consult the community for help.",
			err)
	}

	return nil
}

// TODO: test
func (i *Infrastructure) Compile(ctx *infrastructure.Context) (*infrastructure.CompileResult, error) {
	if err := i.Bindata.CopyDir(ctx.Dir, "data/"+ctx.Infra.Flavor); err != nil {
		return nil, err
	}

	return nil, nil
}

// TODO: impl and test
func (i *Infrastructure) Flavors() []string {
	return nil
}
