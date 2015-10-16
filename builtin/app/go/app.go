package goapp

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/oneline"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=goapp -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	var opts compile.AppOptions
	custom := &customizations{Opts: &opts}
	opts = compile.AppOptions{
		Ctx: ctx,
		FoundationConfig: foundation.Config{
			ServiceName: ctx.Application.Name,
		},
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context: map[string]interface{}{
				"dep_binary_path": fmt.Sprintf("/usr/local/bin/%s", ctx.Application.Name),
				"path": map[string]string{
					"guest_working": fmt.Sprintf(
						"/otto-deps/%s-%s",
						ctx.Application.Name,
						ctx.Appfile.ID),
				},
			},
		},
		Customizations: []*compile.Customization{
			&compile.Customization{
				Type:     "go",
				Callback: custom.processGo,
				Schema: map[string]*schema.FieldSchema{
					"go_version": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "1.5",
						Description: "Go version to install",
					},

					"import_path": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "",
						Description: "Go import path for where to put this in the GOPATH",
					},
				},
			},

			&compile.Customization{
				Type:     "dev-dep",
				Callback: custom.processDevDep,
				Schema: map[string]*schema.FieldSchema{
					"run_command": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "{{ dep_binary_path }}",
						Description: "Command to run this app as a dep",
					},
				},
			},
		},
	}

	return compile.App(&opts)
}

func (a *App) Build(ctx *app.Context) error {
	return fmt.Errorf(strings.TrimSpace(buildErr))
}

func (a *App) Deploy(ctx *app.Context) error {
	return fmt.Errorf(strings.TrimSpace(buildErr))
}

func (a *App) Dev(ctx *app.Context) error {
	// Read the go version, since we use that for our layer
	goVersion, err := oneline.Read(filepath.Join(ctx.Dir, "dev", "go_version"))
	if err != nil {
		return err
	}

	// Setup layers
	layered, err := vagrant.DevLayered(ctx, []*vagrant.Layer{
		&vagrant.Layer{
			ID:          fmt.Sprintf("go%s", goVersion),
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
	return vagrant.DevDep(dst, src, &vagrant.DevDepOptions{
		Dir:    filepath.Join(src.Dir, "dev-dep"),
		Script: "/otto/build.sh",
		Files:  []string{"dev-dep-output"},
	})
}

const devInstructions = `
A development environment has been created for writing a generic Go-based
application. For this development environment, Go is pre-installed. To
work on your project, edit files locally on your own machine. The file changes
will be synced to the development environment.

When you're ready to build your project, run 'otto dev ssh' to enter
the development environment. You'll be placed directly into the working
directory where you can run 'go get' and 'go build' as you normally would.
The GOPATH is already completely setup.
`

const buildErr = `
Build isn't supported yet for Go!

Early versions of Otto are focusing on creating a fantastic development
experience. Because of this, build/deploy are still lacking for many
application types. These will be fixed very soon in upcoming versions of
Otto. Sorry!
`
