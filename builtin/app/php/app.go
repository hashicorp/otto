package phpapp

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

//go:generate go-bindata -pkg=phpapp -nomemcopy -nometadata ./data/...

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
				Type:     "php",
				Callback: custom.processPhp,
				Schema: map[string]*schema.FieldSchema{
					"php_version": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "5.6",
						Description: "PHP version to install",
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
A development environment has been created for writing a PHP app.

Edit files locally on your machine, the file changes will be synced
to the development environment automatically.

To run and view your application, run 'otto dev ssh' to enter the
development environment. You'll be placed directly into the working
directory where you can run "composer", "php", etc.

You can access the environment from this machine using the IP address above.
For example, if you start your app with 'php -S 0.0.0.0:5000', then you can
access it using the above IP at port 5000.
`

const buildErr = `
Build isn't supported yet for PHP!

Early versions of Otto are focusing on creating a fantastic development
experience. Because of this, build/deploy are still lacking for many
application types. These will be fixed very soon in upcoming versions of
Otto. Sorry!
`
