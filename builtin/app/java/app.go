package javaapp

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	javaSP "github.com/hashicorp/otto/builtin/scriptpack/java"
	stdSP "github.com/hashicorp/otto/builtin/scriptpack/stdlib"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/oneline"
	"github.com/hashicorp/otto/helper/packer"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/helper/vagrant"
	"github.com/hashicorp/otto/scriptpack"
)

//go:generate go-bindata -pkg=javaapp -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Meta() (*app.Meta, error) {
	return Meta, nil
}

func (a *App) Implicit(ctx *app.Context) (*appfile.File, error) {
	return nil, nil
}

// Compile ...
func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	var opts compile.AppOptions
	custom := &customizations{Opts: &opts}
	opts = compile.AppOptions{
		Ctx: ctx,
		Result: &app.CompileResult{
			Version: 1,
		},
		FoundationConfig: foundation.Config{
			ServiceName: ctx.Application.Name,
		},
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context:  map[string]interface{}{},
		},
		ScriptPacks: []*scriptpack.ScriptPack{
			&stdSP.ScriptPack,
			&javaSP.ScriptPack,
		},
		Customization: (&compile.Customization{
			Callback: custom.process,
			Schema: map[string]*schema.FieldSchema{
				"java_version": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     "1.8.0_72",
					Description: "Java version installed",
				},
				"gradle_version": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     "2.12",
					Description: "Gradle version to install",
				},
				"maven_version": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     "3.3.9",
					Description: "Maven version to install",
				},
				"scala_version": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     "2.11.7",
					Description: "Scala version to install",
				},
				"sbt_version": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     "0.13.9",
					Description: "sbt version to install",
				},
				"lein_version": &schema.FieldSchema{
					Type:        schema.TypeString,
					Default:     "2.5.3",
					Description: "Leiningen version to install",
				},
			},
		}).Merge(compile.VagrantCustomizations(&opts)),
	}

	return compile.App(&opts)
}

// Build ...
func (a *App) Build(ctx *app.Context) error {
	return packer.Build(ctx, &packer.BuildOptions{
	})
}

// Deploy ...
func (a *App) Deploy(ctx *app.Context) error {
	return terraform.Deploy(&terraform.DeployOptions{
	}).Route(ctx)
}

// Dev ...
func (a *App) Dev(ctx *app.Context) error {
	var layered *vagrant.Layered

	// We only setup a layered environment if we've recompiled since
	// version 0. If we're still at version 0 then we have to use the
	// non-layered dev environment.
	if ctx.CompileResult.Version > 0 {
		// Read the go version, since we use that for our layer
		javaVersion, err := oneline.Read(filepath.Join(ctx.Dir, "dev", "java_version"))
		if err != nil {
			return err
		}

		// Setup layers
		layered, err = vagrant.DevLayered(ctx, []*vagrant.Layer{
			&vagrant.Layer{
				ID:          fmt.Sprintf("java%s", javaVersion),
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

// DevDep ...
func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return vagrant.DevDep(dst, src, &vagrant.DevDepOptions{})
}

const devInstructions = `
A development environment has been created for writing a generic Java based
application using Java as the build system. For this development environment,
Java is pre-installed. To work on your project, edit files locally on your own
machine. The file changes will be synced to the development environment.

When you're ready to build or test your project, run 'otto dev ssh'
to enter the development environment. You'll be placed directly into the
working directory where you can run "gradle init", "gradle build", "mvn clean",
"mvn test" etc.

You can access the environment from this machine using the IP address above.
For example, if your app is running on port 5000, then access it using the
IP above plus that port.
`
