package terraform

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/directory"
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

// Deploy deploys an application using Terraform.
//
// This will verify the infrastruction is created and a build is available,
// and use that information to run Terraform. Any edge cases around Terraform
// failures is handled and state storage is automatic as well.
//
// This function implements app.App.Deploy.
func Deploy(ctx *app.Context, opts *DeployOptions) error {
	// Get the infrastructure state. The infrastructure must be
	// created for us to deploy to it.
	infra, err := ctx.Directory.GetInfra(&directory.Infra{
		Lookup: directory.Lookup{
			Infra: ctx.Appfile.ActiveInfrastructure().Name}})
	if err != nil {
		return err
	}

	if infra == nil || infra.State != directory.InfraStateReady {
		return fmt.Errorf(
			"Infrastructure for this application hasn't been built yet.\n" +
				"The deploy step requires this because the target infrastructure\n" +
				"as well as its final properties can affect the deploy process.\n" +
				"Please run `otto infra` to build the underlying infrastructure,\n" +
				"then run `otto deploy` again.")
	}

	// Construct the variables for Terraform from our queried infra
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

	if !opts.DisableBuild {
		// Get the build information. We must have had a prior build in
		// order to deploy.
		build, err := ctx.Directory.GetBuild(&directory.Build{
			App:         ctx.Tuple.App,
			Infra:       ctx.Tuple.Infra,
			InfraFlavor: ctx.Tuple.InfraFlavor,
		})
		if err != nil {
			return err
		}
		if build == nil {
			return fmt.Errorf(
				"This application hasn't been built yet. Please run `otto build`\n" +
					"first so that the deploy step has an artifact to deploy.")
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
			return fmt.Errorf(
				"Unknown deployment target infrastructure: %s\n\n"+
					"This app currently doesn't know how to deploy to this infrastructure.\n"+
					"Please report this to the project.",
				ctx.Tuple.Infra)
		}
		artifactVars, err := ext(ctx, build, infra)
		if err != nil {
			return err
		}
		for k, v := range artifactVars {
			vars[k] = v
		}
	}

	// Get our old deploy to populate the old state data if we have it.
	// This step is critical to make sure that Terraform remains idempotent
	// and that it handles migrations properly.
	//
	// If we don't have a prior deploy, that is okay, we just create one
	// now (with the DeployStateNew to note that we've never deployed). This
	// gives us the UUID we can use for the state storage.
	deployLookup := &directory.Deploy{
		App:         ctx.Tuple.App,
		Infra:       ctx.Tuple.Infra,
		InfraFlavor: ctx.Tuple.InfraFlavor,
	}
	deploy, err := ctx.Directory.GetDeploy(&directory.Deploy{
		App:         ctx.Tuple.App,
		Infra:       ctx.Tuple.Infra,
		InfraFlavor: ctx.Tuple.InfraFlavor,
	})
	if err != nil {
		return err
	}
	if deploy == nil {
		// If we have no deploy, put in a temporary one
		deploy = deployLookup
		deploy.State = directory.DeployStateNew

		// Write the temporary deploy so we have an ID to use for the state
		if err := ctx.Directory.PutDeploy(deploy); err != nil {
			return err
		}
	}

	// Get the directory
	tfDir := opts.Dir
	if tfDir == "" {
		tfDir = filepath.Join(ctx.Dir, "deploy")
	}

	// Run Terraform!
	tf := &Terraform{
		Dir:       tfDir,
		Ui:        ctx.Ui,
		Variables: vars,
		Directory: ctx.Directory,
		StateId:   deploy.ID,
	}
	if err := tf.Execute("apply"); err != nil {
		return fmt.Errorf(
			"Error running Terraform: %s\n\n"+
				"Terraform usually has helpful error messages. Please read the error\n"+
				"messages above and resolve them. Sometimes simply running `otto deply`\n"+
				"again will work.",
			err)
	}

	return nil
}
