package packer

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/atlas-go/archive"
	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/foundation"
)

type BuildOptions struct {
	// Dir is the directory where Packer will be executed from.
	// If this isn't set, it'll default to "#{ctx.Dir}/build"
	Dir string

	// The path to the template to execute. If this isn't set, it'll
	// default to "#{Dir}/template.json"
	TemplatePath string

	// InfraOutputMap is a map to change the key of an infra output
	// to a different key for a Packer variable. The key of this map
	// is the infra output key, and teh value is the Packer variable name.
	InfraOutputMap map[string]string
}

// Build can be used to build an artifact with Packer and parse the
// artifact out into a Build properly.
//
// This function automatically knows how to parse various built-in
// artifacts of Packer. For the exact functionality of the parse
// functions, see the documentation of the various parse functions.
//
// This function implements the app.App.Build function.
// TODO: Test
func Build(ctx *app.Context, opts *BuildOptions) error {
	project := Project(&ctx.Shared)
	if err := project.InstallIfNeeded(); err != nil {
		return err
	}

	ctx.Ui.Header("Querying infrastructure data for build...")

	// Get the infrastructure, since it needs to be ready for building
	// to occur. We'll copy the outputs and the credentials as variables
	// to Packer.
	infra, err := ctx.Directory.GetInfra(&directory.Infra{
		Lookup: directory.Lookup{
			Infra: ctx.Appfile.ActiveInfrastructure().Name}})
	if err != nil {
		return err
	}

	// If the infra isn't ready then we can't build
	if infra == nil || infra.State != directory.InfraStateReady {
		return fmt.Errorf(
			"Infrastructure for this application hasn't been built yet.\n" +
				"The build step requires this because the target infrastructure\n" +
				"as well as its final properties can affect the build process.\n" +
				"Please run `otto infra` to build the underlying infrastructure,\n" +
				"then run `otto build` again.")
	}

	// Construct the variables for Packer from the infra. We copy them as-is.
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

	// Setup the vars
	if err := foundation.WriteVars(&ctx.Shared); err != nil {
		return fmt.Errorf("Error preparing build: %s", err)
	}

	ctx.Ui.Header("Building deployment archive...")
	slugPath, err := createAppSlug(filepath.Dir(ctx.Appfile.Path))
	if err != nil {
		return err
	}
	vars["slug_path"] = slugPath

	// Start building the resulting build
	build := &directory.Build{
		Lookup: directory.Lookup{
			AppID:       ctx.Appfile.ID,
			Infra:       ctx.Tuple.Infra,
			InfraFlavor: ctx.Tuple.InfraFlavor,
		},

		Artifact: make(map[string]string),
	}

	// Get the paths for Packer execution
	packerDir := opts.Dir
	templatePath := opts.TemplatePath
	if opts.Dir == "" {
		packerDir = filepath.Join(ctx.Dir, "build")
	}
	if opts.TemplatePath == "" {
		templatePath = filepath.Join(packerDir, "template.json")
	}

	ctx.Ui.Header("Building deployment artifact with Packer...")
	ctx.Ui.Message(
		"Raw Packer output will begin streaming in below. Otto\n" +
			"does not create this output. It is mirrored directly from\n" +
			"Packer while the build is being run.\n\n")

	// Build and execute Packer
	p := &Packer{
		Path:      project.Path(),
		Dir:       packerDir,
		Ui:        ctx.Ui,
		Variables: vars,
		Callbacks: map[string]OutputCallback{
			"artifact": ParseArtifactAmazon(build.Artifact),
		},
	}
	if err := p.Execute("build", templatePath); err != nil {
		return err
	}

	// Store the build!
	ctx.Ui.Header("Storing build data in directory...")
	if err := ctx.Directory.PutBuild(build); err != nil {
		return fmt.Errorf(
			"Error storing the build in the directory service: %s\n\n"+
				"Despite the build itself completing successfully, Otto must\n"+
				"also successfully store the results in the directory service\n"+
				"to be able to deploy this build. Please fix the above error and\n"+
				"rebuild.",
			err)
	}

	ctx.Ui.Header("[green]Build success!")
	ctx.Ui.Message(
		"[green]The build was completed successfully and stored within\n" +
			"the directory service, meaning other members of your team\n" +
			"don't need to rebuild this same version and can deploy it\n" +
			"immediately.")

	return nil
}

// ParseArtifactAmazon parses AMIs out of the output.
//
// The map will be populated where the key is the region and the value is
// the AMI ID.
func ParseArtifactAmazon(m map[string]string) OutputCallback {
	return func(o *Output) {
		// We're looking for ID events.
		//
		// Example: 1440649959,amazon-ebs,artifact,0,id,us-east-1:ami-9d66def6
		if len(o.Data) < 3 || o.Data[1] != "id" {
			return
		}

		// TODO: multiple AMIs
		parts := strings.Split(o.Data[2], ":")
		m[parts[0]] = parts[1]
	}
}

// createAppSlug makes an archive of the app with (otto-specific exclusions)
// and yields a path to a tempfile containing that archive
//
// TODO: allow customization of the Exclude patterns
func createAppSlug(path string) (string, error) {
	archive, err := archive.CreateArchive(path, &archive.ArchiveOpts{
		Exclude: []string{".otto", ".vagrant"},
		VCS:     true,
	})
	if err != nil {
		return "", err
	}
	defer archive.Close()

	// Archive is just a reader, and we need it in a file. The below seems
	// fiddly, could there be a better way?
	slug, err := ioutil.TempFile("", "otto-slug-")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(slug, archive)
	cerr := slug.Close()
	if err != nil {
		return "", err
	}
	if cerr != nil {
		return "", err
	}

	return slug.Name(), nil
}
