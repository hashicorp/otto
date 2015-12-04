package phpapp

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/oneline"
	"github.com/hashicorp/otto/helper/packer"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=phpapp -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Meta() (*app.Meta, error) {
	return Meta, nil
}

func (a *App) Implicit(ctx *app.Context) (*appfile.File, error) {
	// For Wordpress we implicitly depend on MySQL
	var result appfile.File
	if ctx.Tuple.App == "wordpress" {
		result.Application = &appfile.Application{
			Dependencies: []*appfile.Dependency{
				&appfile.Dependency{
					Source: "github.com/hashicorp/otto/examples/mysql",
				},
			},
		}
	}

	return &result, nil
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
		Customization: (&compile.Customization{
			Callback: custom.process,
			Schema: map[string]*schema.FieldSchema{
				"php_version": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     "5.6",
					Description: "PHP version to install",
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
	var layered *vagrant.Layered

	// We only setup a layered environment if we've recompiled since
	// version 0. If we're still at version 0 then we have to use the
	// non-layered dev environment.
	if ctx.CompileResult.Version > 0 {
		// Read the go version, since we use that for our layer
		version, err := oneline.Read(filepath.Join(ctx.Dir, "dev", "php_version"))
		if err != nil {
			return err
		}

		// Setup layers
		layered, err = vagrant.DevLayered(ctx, []*vagrant.Layer{
			&vagrant.Layer{
				ID:          fmt.Sprintf("php%s", version),
				Vagrantfile: filepath.Join(ctx.Dir, "dev", "layer-base", "Vagrantfile"),
			},
		})
		if err != nil {
			return err
		}
	}

	instructions := devInstructions
	if ctx.Tuple.App == "wordpress" {
		instructions = devInstructionsWordpress
	}

	// Build the actual development environment
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(instructions),
		Layer:        layered,
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

const devInstructionsWordpress = `
A development environment has been created for working on Wordpress.

To start the web server, SSH into the development environment using
"otto dev ssh" and run "php -S 0.0.0.0:3000". You can then visit Wordpress
using the IP above on port 3000.

MySQL has also automatically been setup. The address for MySQL is
"mysql.service.consul", the username and password is "root".
`

const buildErr = `
Build isn't supported yet for PHP!

Early versions of Otto are focusing on creating a fantastic development
experience. Because of this, build/deploy are still lacking for many
application types. These will be fixed very soon in upcoming versions of
Otto. Sorry!
`
