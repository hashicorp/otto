package scriptpackapp

import (
	"fmt"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=scriptpackapp -nomemcopy -nometadata ./data/...

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
	opts = compile.AppOptions{
		Ctx: ctx,
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
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
	// Build the actual development environment
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
	}).Route(ctx)
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return nil, fmt.Errorf("A ScriptPack can't be a dependency")
}

const devInstructions = `
A development environment has been created for testing this ScriptPack.

Anytime you change the contents of the ScriptPack, you must run
"otto dev scriptpack-rebuild". This will update the contents in the
dev environment.

To run tests, you can use "otto dev scriptpack-test" with the path
to the directory or BATS test file to run. This will be automatically
configured and run within the dev environment.
`

const buildErr = `
Build and deploy aren't supported for ScriptPacks since it doesn't
make a lot of sense.
`
