package nodeapp

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

//go:generate go-bindata -pkg=nodeapp -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

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
		Customizations: []*compile.Customization{
			&compile.Customization{
				Type:     "node",
				Callback: custom.processDev,
				Schema: map[string]*schema.FieldSchema{
					"node_version": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "4.1.0",
						Description: "Node version to install",
					},
				},
			},
		},
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
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
	}).Route(ctx)
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return vagrant.DevDep(dst, src, &vagrant.DevDepOptions{})
}

const devInstructions = `
A development environment has been created for writing a generic Node-based
application. For this development environment, Node is pre-installed. To
work on your project, edit files locally on your own machine. The file changes
will be synced to the development environment.

When you're ready to build or test your project, run 'otto dev ssh'
to enter the development environment. You'll be placed directly into the
working directory where you can run "npm install", "npm run", etc.

You can access the environment from this machine using the IP address above.
For example, if your app is running on port 5000, then access it using the
IP above plus that port.
`
