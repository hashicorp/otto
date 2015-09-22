package terraform

import (
	"fmt"
	"strings"

	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/router"
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
	r := &router.Router{
		Actions: map[string]router.Action{
			"": &router.SimpleAction{
				ExecuteFunc:  i.actionApply,
				SynopsisText: infraApplySyn,
				HelpText:     strings.TrimSpace(infraApplyHelp),
			},
			"destroy": &router.SimpleAction{
				ExecuteFunc:  i.actionDestroy,
				SynopsisText: infraDestroySyn,
				HelpText:     strings.TrimSpace(infraDestroyHelp),
			},
			"info": &router.SimpleAction{
				ExecuteFunc:  i.actionInfo,
				SynopsisText: infraInfoSyn,
				HelpText:     strings.TrimSpace(infraInfoHelp),
			},
		},
	}
	return r.Route(ctx)
}

func (i *Infrastructure) actionDestroy(rctx router.Context) error {
	rctx.UI().Header("Destroying main infrastructure...")
	ctx := rctx.(*infrastructure.Context)
	return i.execute(ctx, "destroy", "-force")
}

func (i *Infrastructure) actionApply(rctx router.Context) error {
	rctx.UI().Header("Building main infrastructure...")
	ctx := rctx.(*infrastructure.Context)
	return i.execute(ctx, "apply")
}

func (i *Infrastructure) actionInfo(rctx router.Context) error {
	ctx := rctx.(*infrastructure.Context)
	project, err := Project(&ctx.Shared)
	if err != nil {
		return err
	}

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
		return fmt.Errorf("Infrastructure not created. Nothing to display.")
	}

	tf := &Terraform{
		Path:      project.Path(),
		Dir:       ctx.Dir,
		Ui:        ctx.Ui,
		Directory: ctx.Directory,
		StateId:   infra.ID,
	}

	// Start the Terraform command
	args := make([]string, len(ctx.ActionArgs)+1)
	args[0] = "output"
	copy(args[1:], ctx.ActionArgs)
	if err := tf.Execute(args...); err != nil {
		return fmt.Errorf("Error running Terraform: %s", err)
	}
	return nil
}

func (i *Infrastructure) execute(ctx *infrastructure.Context, command ...string) error {
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
	err = tf.Execute(command...)
	if err != nil {
		err = fmt.Errorf("Error running Terraform: %s", err)
		infra.State = directory.InfraStatePartial
	}

	ctx.Ui.Header("Terraform execution complete. Saving results...")

	if err == nil {
		if ctx.Action == "destroy" {
			// If we just destroyed successfully, the infra is now empty.
			infra.State = directory.InfraStateInvalid
			infra.Outputs = map[string]string{}
		} else {
			// If an apply was successful, populate the state and outputs.
			infra.State = directory.InfraStateReady
			infra.Outputs, err = tf.Outputs()
			if err != nil {
				err = fmt.Errorf("Error reading Terraform outputs: %s", err)
				infra.State = directory.InfraStatePartial
			}
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

// Synopsis text for actions
const (
	infraApplySyn   = "Create or update infrastructure resources for this application"
	infraDestroySyn = "Destroy infrastructure resources for this application"
	infraInfoSyn    = "Display information about this application's infrastructure"
)

// Help text for actions
const infraApplyHelp = `
Usage: otto infra

  Creates infrastructure for your application.

  This command will create all the resource required to serve as an
  infrastructure for your application.
`

const infraDestroyHelp = `
Usage: otto infra destroy [-force]

  Destroys all infrastructure resources.

  This command will remove any previously-created infrastructure resources.
  Note that any apps with resources deployed into this infrastructure will need
  to have 'otto deploy destroy' run before this command will succeed.

	Otto will ask for confirmation to protect against an accidental destroy. You
	can provide the -force flag to skip this check.
`

const infraInfoHelp = `
Usage: otto infra info [NAME]

  Displays information about this application's infrastructure.

  This command will show any variables the infrastructure has specified as
  outputs. If no NAME is specified, all outputs will be listed. If NAME is
  specified, just the contents of that output will be printed.
`
