package rubyapp

import (
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/packer"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=rubyapp -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	return compile.App(ctx, &compile.AppOptions{
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context:  map[string]interface{}{},
		},
		Customizations: []*compile.Customization{
			&compile.Customization{
				Type:     "dev",
				Callback: processCustomDev,
				Schema: map[string]*schema.FieldSchema{
					"ruby_version": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "2.2",
						Description: "Ruby version to install",
					},
				},
			},
		},
	})
}

func (a *App) Build(ctx *app.Context) error {
	return packer.Build(ctx, &packer.BuildOptions{
		InfraOutputMap: map[string]string{
			"region": "aws_region",
		},
	})
}

func (a *App) Deploy(ctx *app.Context) error {
	return terraform.Deploy(ctx, &terraform.DeployOptions{
		InfraOutputMap: map[string]string{
			"region":         "aws_region",
			"subnet-private": "private_subnet_id",
			"subnet-public":  "public_subnet_id",
		},
	})
}

func (a *App) Dev(ctx *app.Context) error {
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
	}).Route(ctx)
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return vagrant.DevDep(dst, src, &vagrant.DevDepOptions{})
}

func processCustomDev(d *schema.FieldData) (*compile.CustomizationResult, error) {
	return &compile.CustomizationResult{
		TemplateContext: map[string]interface{}{
			"dev_ruby_version": d.Get("ruby_version"),
		},
	}, nil
}

const devInstructions = `
A development environment has been created for writing a generic Ruby-based
application. For this development environment, Ruby is pre-installed. To
work on your project, edit files locally on your own machine. The file changes
will be synced to the development environment.

When you're ready to build your project, run 'otto dev ssh' to enter
the development environment. You'll be placed directly into the working
directory where you can run 'bundle install' and 'ruby' as you normally would.
`
