package goapp

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/packer"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=goapp -nomemcopy ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	data := &bindata.Data{
		Asset:    Asset,
		AssetDir: AssetDir,
		Context: map[string]interface{}{
			"name":          ctx.Appfile.Application.Name,
			"dev_fragments": ctx.DevDepFragments,
			"path": map[string]string{
				"cache":    ctx.CacheDir,
				"compiled": ctx.Dir,
				"working":  filepath.Dir(ctx.Appfile.Path),
			},
		},
	}

	// Copy all the common files
	if err := data.CopyDir(ctx.Dir, "data/common"); err != nil {
		return nil, err
	}

	// Copy the infrastructure specific files
	prefix := fmt.Sprintf("data/%s-%s", ctx.Tuple.Infra, ctx.Tuple.InfraFlavor)
	if err := data.CopyDir(ctx.Dir, prefix); err != nil {
		return nil, err
	}

	return &app.CompileResult{
		DevDepFragmentPath: filepath.Join(ctx.Dir, "dev-dep/build/Vagrantfile.fragment"),
	}, nil
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
	return vagrant.Dev(ctx, &vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
	})
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return vagrant.DevDep(dst, src, &vagrant.DevDepOptions{
		Dir:    filepath.Join(src.Dir, "dev-dep/build"),
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
