package pythonapp

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	pythonSP "github.com/hashicorp/otto/builtin/scriptpack/python"
	stdSP "github.com/hashicorp/otto/builtin/scriptpack/stdlib"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/oneline"
	"github.com/hashicorp/otto/helper/packer"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/helper/vagrant"
	"github.com/hashicorp/otto/scriptpack"
)

//go:generate go-bindata -pkg=pythonapp -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Meta() (*app.Meta, error) {
	return Meta, nil
}

func (a *App) Implicit(ctx *app.Context) (*appfile.File, error) {
	return nil, nil
}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	var opts compile.AppOptions
	custom := &customizations{Opts: &opts}
	opts = compile.AppOptions{
		Ctx: ctx,
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context:  map[string]interface{}{},
		},
		ScriptPacks: []*scriptpack.ScriptPack{
			&stdSP.ScriptPack,
			&pythonSP.ScriptPack,
		},
		Customization: (&compile.Customization{
			Callback: custom.process,
			Schema: map[string]*schema.FieldSchema{
				"python_version": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     "2.7",
					Description: "Python version to install",
				},
				"python_entrypoint": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     fmt.Sprintf("%s:app", ctx.Appfile.Application.Name),
					Description: "WSGI entry point",
				},
			},
		}).Merge(compile.VagrantCustomizations(&opts)),
	}

	return compile.App(&opts)
}

func (a *App) Build(ctx *app.Context) error {
	return packer.Build(ctx, &packer.BuildOptions{
		InfraOutputMap: map[string]string{
			"region": "aws_region",
		},
	})
}

func (a *App) Deploy(ctx *app.Context) error {
	return terraform.Deploy(&terraform.DeployOptions{
		InfraOutputMap: map[string]string{
			"region":         "aws_region",
			"subnet-private": "private_subnet_id",
			"subnet-public":  "public_subnet_id",
		},
	}).Route(ctx)
}

func (a *App) Dev(ctx *app.Context) error {
	// Read the go version, since we use that for our layer
	version, err := oneline.Read(filepath.Join(ctx.Dir, "dev", "python_version"))
	if err != nil {
		return err
	}

	// Setup layers
	layered, err := vagrant.DevLayered(ctx, []*vagrant.Layer{
		&vagrant.Layer{
			ID:          fmt.Sprintf("python%s", version),
			Vagrantfile: filepath.Join(ctx.Dir, "dev", "layer-base", "Vagrantfile"),
		},
	})
	if err != nil {
		return err
	}

	// Build the actual development environment
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
		Layer:        layered,
	}).Route(ctx)
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return vagrant.DevDep(dst, src, &vagrant.DevDepOptions{})
}

const devInstructions = `
A development environment has been created for writing a generic
Python-based app.

Python is pre-installed. To work on your project, edit files locally on your
own machine. The file changes will be synced to the development environment.

When you're ready to build your project, run 'otto dev ssh' to enter
the development environment. You'll be placed directly into the working
directory with a virtualenv setup where you can run 'pip' and 'python'
as you normally would.

You can access any running web application using the IP above.
`
