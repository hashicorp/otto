package terraform

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/helper/router"
)

type DeployOptions struct {
	// Dir is the directory where Terraform is run. If this isn't set, it'll
	// default to "#{ctx.Dir}/deploy".
	Dir string

	// DisableBuild, if true, will not load a build associated with this
	// appfile and attempt to extract the artifact from it. In this case,
	// AritfactExtractors is also useless.
	DisableBuild bool

	// ArtifactExtractors is a mapping of artifact extractors. The
	// built-in artifact extractors will populate this if a key isn't set.
	ArtifactExtractors map[string]DeployArtifactExtractor

	// InfraOutputMap is a map to change the key of an infra output
	// to a different key for a Terraform variable. The key of this map
	// is the infra output key, and teh value is the Terraform variable name.
	InfraOutputMap map[string]string
}

// Deploy can be used as an implementation of app.App.Deploy to handle calling
// out to terraform w/ the configured config to get an app deployed to an
// infrastructure.
//
// This will verify the infrastructure is created and a build is available,
// and use that information to run Terraform. Any edge cases around Terraform
// failures is handled and state storage is automatic as well.
//
// This function implements app.App.Deploy.
func Deploy(opts *DeployOptions) *router.Router {
	return &router.Router{
		Actions: map[string]router.Action{
			"": &router.SimpleAction{
				ExecuteFunc:  opts.actionDeploy,
				SynopsisText: actionDeploySyn,
				HelpText:     strings.TrimSpace(actionDeployHelp),
			},
			"destroy": &router.SimpleAction{
				ExecuteFunc:  opts.actionDestroy,
				SynopsisText: actionDestroySyn,
				HelpText:     strings.TrimSpace(actionDestroyHelp),
			},
			"info": &router.SimpleAction{
				ExecuteFunc:  opts.actionInfo,
				SynopsisText: actionInfoSyn,
				HelpText:     strings.TrimSpace(actionInfoHelp),
			},
		},
	}
}

func (opts *DeployOptions) actionDeploy(rctx router.Context) error {
	ctx := rctx.(*app.Context)
	project, err := Project(&ctx.Shared)
	if err != nil {
		return err
	}
	vars := make(map[string]string)

	infra, infraVars, err := opts.lookupInfraVars(ctx)
	if err != nil {
		return err
	}
	if infra == nil {
		return fmt.Errorf(
			"Infrastructure for this application hasn't been built yet.\n" +
				"The deploy step requires this because the target infrastructure\n" +
				"as well as its final properties can affect the deploy process.\n" +
				"Please run `otto infra` to build the underlying infrastructure,\n" +
				"then run `otto deploy` again.")
	}
	for k, v := range infraVars {
		vars[k] = v
	}

	if !opts.DisableBuild {
		buildVars, err := opts.lookupBuildVars(ctx, infra)
		if err != nil {
			return err
		}
		if buildVars == nil {
			return fmt.Errorf(
				"This application hasn't been built yet. Please run `otto build`\n" +
					"first so that the deploy step has an artifact to deploy.")
		}
		for k, v := range buildVars {
			vars[k] = v
		}
	}

	// Setup the vars
	if err := foundation.WriteVars(&ctx.Shared); err != nil {
		return fmt.Errorf("Error preparing deploy: %s", err)
	}

	// Get our old deploy to populate the old state data if we have it.
	// This step is critical to make sure that Terraform remains idempotent
	// and that it handles migrations properly.
	deploy, err := opts.lookupDeploy(ctx)
	if err != nil {
		return err
	}

	// Run Terraform!
	tf := &Terraform{
		Path:      project.Path(),
		Dir:       opts.tfDir(ctx),
		Ui:        ctx.Ui,
		Variables: vars,
		Directory: ctx.Directory,
		StateId:   deploy.ID,
	}
	if err := tf.Execute("apply"); err != nil {
		deploy.MarkFailed()
		if putErr := ctx.Directory.PutDeploy(deploy); putErr != nil {
			return fmt.Errorf("The deploy failed with err: %s\n\n"+
				"And then there was an error storing it in the directory: %s\n"+
				"This second error is a bug and should be reported.", err, putErr)
		}

		return terraformError(err)
	}

	deploy.MarkSuccessful()
	if err := ctx.Directory.PutDeploy(deploy); err != nil {
		return err
	}
	return nil
}

func (opts *DeployOptions) actionDestroy(rctx router.Context) error {
	ctx := rctx.(*app.Context)
	project, err := Project(&ctx.Shared)
	if err != nil {
		return err
	}
	vars := make(map[string]string)

	infra, infraVars, err := opts.lookupInfraVars(ctx)
	if err != nil {
		return err
	}
	if infra == nil {
		return fmt.Errorf(
			"Infrastructure for this application hasn't been built yet.\n" +
				"Nothing to destroy.")
	}
	for k, v := range infraVars {
		vars[k] = v
	}

	if !opts.DisableBuild {
		buildVars, err := opts.lookupBuildVars(ctx, infra)
		if err != nil {
			return err
		}
		if buildVars == nil {
			return fmt.Errorf(
				"This application hasn't been built yet. Nothing to destroy.")
		}
		for k, v := range buildVars {
			vars[k] = v
		}
	}

	deploy, err := opts.lookupDeploy(ctx)
	if err != nil {
		return err
	}
	if deploy.IsNew() {
		return fmt.Errorf(
			"This application hasn't been deployed yet. Nothing to destroy.")
	}

	// Get the directory
	// Run Terraform!
	tf := &Terraform{
		Path:      project.Path(),
		Dir:       opts.tfDir(ctx),
		Ui:        ctx.Ui,
		Variables: vars,
		Directory: ctx.Directory,
		StateId:   deploy.ID,
	}
	if err := tf.Execute("destroy"); err != nil {
		deploy.MarkFailed()
		if putErr := ctx.Directory.PutDeploy(deploy); putErr != nil {
			return fmt.Errorf("The destroy failed with err: %s\n\n"+
				"And then there was an error storing it in the directory: %s\n"+
				"This second error is a bug and should be reported.", err, putErr)
		}

		return terraformError(err)
	}

	deploy.MarkGone()
	if err := ctx.Directory.PutDeploy(deploy); err != nil {
		return err
	}

	return nil
}

func (opts *DeployOptions) actionInfo(rctx router.Context) error {
	ctx := rctx.(*app.Context)
	project, err := Project(&ctx.Shared)
	if err != nil {
		return err
	}

	deploy, err := opts.lookupDeploy(ctx)
	if err != nil {
		return err
	}
	if deploy.IsNew() {
		return fmt.Errorf(
			"This application hasn't been deployed yet. Nothing to show.")
	}

	// Get the directory
	// Run Terraform!
	tf := &Terraform{
		Path:      project.Path(),
		Dir:       opts.tfDir(ctx),
		Ui:        ctx.Ui,
		Directory: ctx.Directory,
		StateId:   deploy.ID,
	}
	args := make([]string, len(ctx.ActionArgs)+1)
	args[0] = "output"
	copy(args[1:], ctx.ActionArgs)
	if err := tf.Execute(args...); err != nil {
		return terraformError(err)
	}

	return nil
}

// lookupInfraVars collects information about the result of `otto infra` and
// yields a set of variables that can be used by the deploy to reference
// resources in the infrastructure. It returns `nil` if the infrastructure has
// not been created successfully yet.
func (opts *DeployOptions) lookupInfraVars(
	ctx *app.Context) (*directory.Infra, map[string]string, error) {
	infra, err := ctx.Directory.GetInfra(&directory.Infra{
		Lookup: directory.Lookup{
			Infra: ctx.Appfile.ActiveInfrastructure().Name}})
	if err != nil {
		return nil, nil, err
	}

	if !infra.IsReady() {
		return nil, nil, nil
	}

	vars := make(map[string]string)
	for k, v := range infra.Outputs {
		if opts.InfraOutputMap != nil {
			if nk, ok := opts.InfraOutputMap[k]; ok {
				k = nk
			}
		}
		vars[k] = v
	}
	for k, v := range ctx.InfraCreds {
		vars[k] = v
	}
	return infra, vars, nil
}

// lookupBuildVars collects information about the result of `otto build` and
// yields a set of variables that can be used by the deploy to reference the
// built artifact. It returns nil if `otto build` has not yet been run.
func (opts *DeployOptions) lookupBuildVars(
	ctx *app.Context, infra *directory.Infra) (map[string]string, error) {
	build, err := ctx.Directory.GetBuild(&directory.Build{
		Lookup: directory.Lookup{
			AppID:       ctx.Appfile.ID,
			Infra:       ctx.Tuple.Infra,
			InfraFlavor: ctx.Tuple.InfraFlavor,
		},
	})
	if err != nil {
		return nil, err
	}
	if build == nil {
		return nil, nil
	}

	// Extract the artifact from the build. We do this based on the
	// infrastructure type.
	if opts.ArtifactExtractors == nil {
		opts.ArtifactExtractors = make(map[string]DeployArtifactExtractor)
	}
	for k, v := range deployArtifactExtractors {
		if _, ok := opts.ArtifactExtractors[k]; !ok {
			opts.ArtifactExtractors[k] = v
		}
	}
	ext, ok := opts.ArtifactExtractors[ctx.Tuple.Infra]
	if !ok {
		return nil, fmt.Errorf(
			"Unknown deployment target infrastructure: %s\n\n"+
				"This app currently doesn't know how to deploy to this infrastructure.\n"+
				"Please report this to the project.",
			ctx.Tuple.Infra)
	}
	return ext(ctx, build, infra)
}

// lookupDeploy returns any previously deploy made by Otto so we have the state
// necessary to update it.
//
// If we don't have a prior deploy, that is okay, we just create one
// now (with the DeployStateNew to note that we've never deployed). This
// gives us the UUID we can use for the state storage.
func (opts *DeployOptions) lookupDeploy(
	ctx *app.Context) (*directory.Deploy, error) {
	deployLookup := directory.Lookup{
		AppID:       ctx.Appfile.ID,
		Infra:       ctx.Tuple.Infra,
		InfraFlavor: ctx.Tuple.InfraFlavor,
	}
	deploy, err := ctx.Directory.GetDeploy(&directory.Deploy{Lookup: deployLookup})
	if err != nil {
		return nil, err
	}

	if deploy == nil {
		// If we have no deploy, put in a temporary one
		deploy = &directory.Deploy{Lookup: deployLookup}
		deploy.State = directory.DeployStateNew

		// Write the temporary deploy so we have an ID to use for the state
		if err := ctx.Directory.PutDeploy(deploy); err != nil {
			return nil, err
		}
	}

	return deploy, nil
}

// tfDir returns the appropriate terraform working dir
func (opts *DeployOptions) tfDir(ctx *app.Context) string {
	tfDir := opts.Dir
	if tfDir == "" {
		tfDir = filepath.Join(ctx.Dir, "deploy")
	}
	return tfDir
}

// terraformError wraps an error from Terraform in a friendlier message.
func terraformError(err error) error {
	return fmt.Errorf(
		"Error running Terraform: %s\n\n"+
			"Terraform usually has helpful error messages. Please read the error\n"+
			"messages above and resolve them. Sometimes simply running `otto deploy`\n"+
			"again will work.",
		err)
}

// Synopsis text for actions
const (
	actionDeploySyn  = "Deploy the latest built artifact into your infrastructure"
	actionDestroySyn = "Destroy all deployed resources for this application"
	actionInfoSyn    = "Display information about this application's deploy"
)

// Help text for actions
const actionDeployHelp = `
Usage: otto deploy

  Deploys a built artifact into your infrastructure.

  This command will take the latest built artifact and deploy it into your
  infrastructure. Otto will create or replace any necessary resources required
  to run your app.
`

const actionDestroyHelp = `
Usage: otto deploy destroy

  Destroys any deployed resources associated with this application.

  This command will remove any previously-deployed resources from your
  infrastructure. This must be run for all of apps in an infrastructure before
  'otto infra destroy' will work.
`

const actionInfoHelp = `
Usage: otto deploy info [NAME]

  Displays information about this application's deploy.

	This command will show any variables the deploy has specified as outputs. If
	no NAME is specified, all outputs will be listed. If NAME is specified, just
	the contents of that output will be printed.
`
