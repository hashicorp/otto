package scriptpackapp

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/builtin/app/go"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	execHelper "github.com/hashicorp/otto/helper/exec"
	"github.com/hashicorp/otto/helper/router"
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
	// Get the import path for this
	path, err := goapp.DetectImportPath(ctx)
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, fmt.Errorf(
			"Your ScriptPack development folder must be within your GOPATH like\n" +
				"a standard Go project. This is required for the dev environment\n" +
				"to function properly. Please put this folder in a proper GOPATH\n" +
				"location.")
	}

	var opts compile.AppOptions
	opts = compile.AppOptions{
		Ctx: ctx,
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context: map[string]interface{}{
				"working_gopath": path,
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
	r := vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
	})

	// Add our customer actions
	r.Actions["scriptpack-rebuild"] = &router.SimpleAction{
		ExecuteFunc:  a.actionRebuild,
		SynopsisText: actionRebuildSyn,
		HelpText:     strings.TrimSpace(actionRebuildHelp),
	}

	return r.Route(ctx)
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return nil, fmt.Errorf("A ScriptPack can't be a dependency")
}

func (a *App) actionRebuild(rctx router.Context) error {
	ctx := rctx.(*app.Context)

	// Get the path to the rebuild binary
	path := filepath.Join(ctx.Dir, "rebuild", "rebuild.go")

	// Run it to regenerate the contents
	cmd := exec.Command("go", "run", path)
	if err := execHelper.Run(ctx.Ui, cmd); err != nil {
		return err
	}

	// Success
	ctx.Ui.Raw("ScriptPack data rebuilt!\n")
	return nil
}

const devInstructions = `
A development environment has been created for testing this ScriptPack.

Anytime you change the contents of the ScriptPack, you must run
"otto dev scriptpack-rebuild". This will update the contents in the
dev environment. This is a very fast operation.

To run tests, you can use "otto dev scriptpack-test" with the path
to the directory or BATS test file to run. This will be automatically
configured and run within the dev environment.
`

const buildErr = `
Build and deploy aren't supported for ScriptPacks since it doesn't
make a lot of sense.
`

const (
	actionRebuildSyn = "Rebuild ScriptPack output for dev and test"
)

const actionRebuildHelp = `
Usage: otto dev scriptpack-rebuild

  Rebuilds the ScriptPack files and dependencies into a single directory
  for dev and test within the development environment.

  This command must be run before running tests after making any changes
  to the ScriptPack.

`
