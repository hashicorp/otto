package rubyapp

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	rubySP "github.com/hashicorp/otto/builtin/scriptpack/ruby"
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

//go:generate go-bindata -pkg=rubyapp -nomemcopy -nometadata ./data/...

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
		Result: &app.CompileResult{
			Version: 1,
		},
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context:  map[string]interface{}{},
		},
		ScriptPacks: []*scriptpack.ScriptPack{
			&stdSP.ScriptPack,
			&rubySP.ScriptPack,
		},
		Customization: (&compile.Customization{
			Callback: custom.process,
			Schema: map[string]*schema.FieldSchema{
				"ruby_version": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     "detect",
					Description: "Ruby version to install",
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
			"vpc_id": "aws_vpc_id",
			"subnet_public":  "aws_subnet_id",
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
	var layered *vagrant.Layered

	// We only setup a layered environment if we've recompiled since
	// version 0. If we're still at version 0 then we have to use the
	// non-layered dev environment.
	if ctx.CompileResult.Version > 0 {
		// Read the go version, since we use that for our layer
		version, err := oneline.Read(filepath.Join(ctx.Dir, "dev", "ruby_version"))
		if err != nil {
			return err
		}

		// Setup layers
		layered, err = vagrant.DevLayered(ctx, []*vagrant.Layer{
			&vagrant.Layer{
				ID:          fmt.Sprintf("ruby%s", version),
				Vagrantfile: filepath.Join(ctx.Dir, "dev", "layer-base", "Vagrantfile"),
			},
		})
		if err != nil {
			return err
		}
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
Ruby-based app.

Ruby is pre-installed. To work on your project, edit files locally on your
own machine. The file changes will be synced to the development environment.

When you're ready to build your project, run 'otto dev ssh' to enter
the development environment. You'll be placed directly into the working
directory where you can run 'bundle' and 'ruby' as you normally would.

You can access any running web application using the IP above.
`
