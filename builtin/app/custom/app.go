package custom

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/oneline"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=custom -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	fragmentPath := filepath.Join(ctx.Dir, "dev-dep", "Vagrantfile.fragment")
	return compile.App(ctx, &compile.AppOptions{
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context: map[string]interface{}{
				"fragment_path": fragmentPath,
			},
		},
		Customizations: []*compile.Customization{
			&compile.Customization{
				Type:     "dev-dep",
				Callback: processCustomDevDep,
				Schema: map[string]*schema.FieldSchema{
					"vagrantfile": &schema.FieldSchema{
						Type:        schema.TypeString,
						Description: "Path to Vagrantfile template",
					},
				},
			},

			&compile.Customization{
				Type:     "deploy",
				Callback: processCustomDeploy,
				Schema: map[string]*schema.FieldSchema{
					"terraform": &schema.FieldSchema{
						Type:        schema.TypeString,
						Description: "Path to a Terraform module",
					},
				},
			},
		},
	})
}

func (a *App) Build(ctx *app.Context) error {
	return nil
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

	// But if we did, then deploy using Terraform
	return terraform.Deploy(ctx, &terraform.DeployOptions{
		Dir: tfdir,
	})
}

func (a *App) Dev(ctx *app.Context) error {
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
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

func processCustomDeploy(d *schema.FieldData) (*compile.CustomizationResult, error) {
	tf, ok := d.GetOk("terraform")
	if !ok {
		return nil, nil
	}

	return &compile.CustomizationResult{
		Callback: compileCustomDeploy(d),
		TemplateContext: map[string]interface{}{
			"deploy_terraform_path": tf.(string),
		},
	}, nil
}

func processCustomDevDep(d *schema.FieldData) (*compile.CustomizationResult, error) {
	if _, ok := d.GetOk("vagrantfile"); !ok {
		return nil, nil
	}

	return &compile.CustomizationResult{
		Callback: compileDev(d),
	}, nil
}

func compileCustomDeploy(d *schema.FieldData) compile.CompileCallback {
	return func(ctx *app.Context, data *bindata.Data) error {
		return data.RenderAsset(
			filepath.Join(ctx.Dir, "deploy", "terraform_path"),
			"data/sentinels/terraform_path.tpl")
	}
}

func compileDev(d *schema.FieldData) compile.CompileCallback {
	vf := d.Get("vagrantfile").(string)

	return func(ctx *app.Context, data *bindata.Data) error {
		fragment := data.Context["fragment_path"].(string)
		if err := data.RenderReal(fragment, vf); err != nil {
			return err
		}

		return data.RenderAsset(
			filepath.Join(ctx.Dir, "dev", "Vagrantfile"),
			"data/dev/Vagrantfile.tpl")
	}
}

const devInstructions = `
A development has been created.

Note that this development environment is just an example of what a
consumer of this application might see as a development dependency.
"Custom" types are not meant to be mutably developed like normal
applications.
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
        vagrant = "path/to/Vagrantfile"
    }

`
