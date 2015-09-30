package custom

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/oneline"
	"github.com/hashicorp/otto/helper/packer"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=custom -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	fragmentPath := filepath.Join(ctx.Dir, "dev-dep", "Vagrantfile.fragment")

	var opts compile.AppOptions
	custom := &customizations{Opts: &opts}
	opts = compile.AppOptions{
		Ctx: ctx,
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context: map[string]interface{}{
				"fragment_path": fragmentPath,
			},
		},
		FoundationConfig: foundation.Config{
			ServiceName: ctx.Application.Name,
		},
		Customizations: []*compile.Customization{
			&compile.Customization{
				Type:     "dev",
				Callback: custom.processDev,
				Schema: map[string]*schema.FieldSchema{
					"vagrantfile": &schema.FieldSchema{
						Type:        schema.TypeString,
						Description: "Path to Vagrantfile",
					},
				},
			},

			&compile.Customization{
				Type:     "dev-dep",
				Callback: custom.processDevDep,
				Schema: map[string]*schema.FieldSchema{
					"vagrantfile": &schema.FieldSchema{
						Type:        schema.TypeString,
						Description: "Path to Vagrantfile template",
					},
				},
			},

			&compile.Customization{
				Type:     "build",
				Callback: custom.processBuild,
				Schema: map[string]*schema.FieldSchema{
					"packer": &schema.FieldSchema{
						Type:        schema.TypeString,
						Description: "Path to Packer template",
					},
				},
			},

			&compile.Customization{
				Type:     "deploy",
				Callback: custom.processDeploy,
				Schema: map[string]*schema.FieldSchema{
					"terraform": &schema.FieldSchema{
						Type:        schema.TypeString,
						Description: "Path to a Terraform module",
					},
				},
			},
		},
	}

	return compile.App(&opts)
}

func (a *App) Build(ctx *app.Context) error {
	// Determine if we set a Packer path.
	path := filepath.Join(ctx.Dir, "build", "packer_path")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return errors.New(strings.TrimSpace(errPackerNotSet))
		}

		return err
	}

	// Read the actual Packer dir
	packerPath, err := oneline.Read(path)
	if err != nil {
		return fmt.Errorf(
			"Error reading the Packer template path: %s\n\n"+
				"An Otto recompile with `otto compile` usually fixes this.",
			err)
	}

	return packer.Build(ctx, &packer.BuildOptions{
		Dir:          filepath.Dir(packerPath),
		TemplatePath: packerPath,
	})
}

func (a *App) Deploy(ctx *app.Context) error {
	// Determine if we set a Terraform path. If we didn't, then
	// tell the user we can't deploy.
	path := filepath.Join(ctx.Dir, "deploy", "terraform_path")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return errors.New(strings.TrimSpace(errTerraformNotSet))
		}

		return err
	}

	// Read the actual TF dir
	tfdir, err := oneline.Read(path)
	if err != nil {
		return fmt.Errorf(
			"Error reading the Terraform module directory: %s\n\n"+
				"An Otto recompile with `otto compile` usually fixes this.",
			err)
	}

	// Determine if we set a Packer path. If we didn't, we disable
	// the build loading for deploys.
	disableBuild := true
	path = filepath.Join(ctx.Dir, "build", "packer_path")
	if _, err := os.Stat(path); err == nil {
		disableBuild = false
	}

	// But if we did, then deploy using Terraform
	return terraform.Deploy(&terraform.DeployOptions{
		Dir:          tfdir,
		DisableBuild: disableBuild,
	}).Route(ctx)
}

func (a *App) Dev(ctx *app.Context) error {
	// Determine if we have a Vagrant path set...
	instructions := devInstructionsCustom
	path := filepath.Join(ctx.Dir, "dev", "vagrant_path")
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		path = ""
	}
	if path != "" {
		var err error
		path, err = oneline.Read(path)
		if err != nil {
			return fmt.Errorf(
				"Error reading the Vagrant directory: %s\n\n"+
					"An Otto recompile with `otto compile` usually fixes this.",
				err)
		}
	}

	if path == "" {
		// Determine if we have our own Vagrantfile
		path = filepath.Join(ctx.Dir, "dev", "Vagrantfile")
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				return errors.New(strings.TrimSpace(errVagrantNotSet))
			}

			return err
		}

		instructions = devInstructionsDevDep
	}

	return vagrant.Dev(&vagrant.DevOptions{
		Dir:          filepath.Dir(path),
		Instructions: strings.TrimSpace(instructions),
	}).Route(ctx)
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	// Determine if we have a Vagrantfile. This is a sentinel that
	// we set this setting.
	path := filepath.Join(src.Dir, "dev", "Vagrantfile")
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New(strings.TrimSpace(errVagrantNotSet))
		}

		return nil, err
	}

	// We purposely return nil here. We don't need to do anything. It
	// is all already setup from the compilation step.
	return nil, nil
}

const devInstructionsCustom = `
Vagrant was executed in the directory specified by the "dev"
customization.
`

const devInstructionsDevDep = `
A development has been created.

Note that this development environment is just an example of what a
consumer of this application might see as a development dependency.
"Custom" types are not meant to be mutably developed like normal
applications.
`

const errPackerNotSet = `
Otto can't build this application because the "packer" setting
isn't set in the "build" customization.

For the "custom" application type, the "build" customization must
set the "packer" setting to point to a Packer template to execute.
Otto will execute this for the build.

Example:

    customization "build" {
        packer = "path/to/template.json"
    }

`

const errTerraformNotSet = `
Otto can't deploy this application because the "terraform" setting
hasn't been set in the "deploy" customization.

For the "custom" application type, the "deploy" customization must
set the "terraform" setting to point to a Terraform module to execute.
Otto will execute this module for the deploy.

Example:

    customization "deploy" {
        terraform = "path/to/module"
    }

`

const errVagrantNotSet = `
Otto can't build a development environment for this because the
"vagrant" setting hasn't been set in the "dev" or "dev-dep" customization.

For the "custom" application type, customizations must be used to
tell Otto what to do. For the dev command, Otto requires either
the "dev" or "dev-dep" customization to be set with the "vagrant" setting.
The "vagrant" setting depends on which customization is being set. Please
refer to the documentation for more details.

Example:

    customization "dev" {
        vagrantfile = "path/to/Vagrantfile"
    }

`
