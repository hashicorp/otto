package goapp

import (
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/packer"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=goapp -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	// Go is really finicky about the GOPATH. To help make the dev
	// environment and build environment more correct, we attempt to
	// detect the GOPATH automatically.
	//
	// We use this GOPATH for example in Vagrant to setup the synced
	// folder directly into the GOPATH properly. Magic!
	ctx.Ui.Header("Detecting application import path for GOPATH...")
	gopathPath, err := detectImportPath(ctx)
	if err != nil {
		return nil, err
	}

	folderPath := "/vagrant"
	if gopathPath != "" {
		folderPath = "/opt/gopath/src/" + gopathPath
	}

	var opts compile.AppOptions
	custom := &customizations{Opts: &opts}
	opts = compile.AppOptions{
		Ctx: ctx,
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context: map[string]interface{}{
				"import_path":        gopathPath,
				"shared_folder_path": folderPath,
			},
		},
		Customizations: []*compile.Customization{
			&compile.Customization{
				Type:     "dev",
				Callback: custom.processDev,
				Schema: map[string]*schema.FieldSchema{
					"go_version": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "1.5",
						Description: "Go version to install",
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
	return terraform.Deploy(ctx, &terraform.DeployOptions{})
}

func (a *App) Dev(ctx *app.Context) error {
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
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
