package gradleapp

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

//go:generate go-bindata -pkg=gradleapp -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

// Compile ...
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
				Type:     "gradle",
				Callback: custom.processDev,
				Schema: map[string]*schema.FieldSchema{
					"gradle_version": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "2.7",
						Description: "Gradle version to install",
					},
				},
			},
		},
	}

	return compile.App(&opts)
}

// Build ...
func (a *App) Build(ctx *app.Context) error {
	return packer.Build(ctx, &packer.BuildOptions{
		InfraOutputMap: map[string]string{
			"region": "aws_region",
		},
	})
}

// Deploy ...
func (a *App) Deploy(ctx *app.Context) error {
	return terraform.Deploy(&terraform.DeployOptions{
		InfraOutputMap: map[string]string{
			"region":         "aws_region",
			"subnet-private": "private_subnet_id",
			"subnet-public":  "public_subnet_id",
		},
	}).Route(ctx)
}

// Dev ...
func (a *App) Dev(ctx *app.Context) error {
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
	}).Route(ctx)
}

// DevDep ...
func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return vagrant.DevDep(dst, src, &vagrant.DevDepOptions{})
}

const devInstructions = `
A development environment has been created for writing a generic Gradle-based
application. For this development environment, Gradle is pre-installed. To
work on your project, edit files locally on your own machine. The file changes
will be synced to the development environment.

When you're ready to build or test your project, run 'otto dev ssh'
to enter the development environment. You'll be placed directly into the
working directory where you can run "npm install", "npm run", etc.

You can access the environment from this machine using the IP address above.
For example, if your app is running on port 5000, then access it using the
IP above plus that port.
`
