package scriptpackapp

import (
	"encoding/json"
	"fmt"
	"os"
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

	r.Actions["scriptpack-test"] = &router.SimpleAction{
		ExecuteFunc:  a.actionTest,
		SynopsisText: actionTestSyn,
		HelpText:     strings.TrimSpace(actionTestHelp),
	}

	return r.Route(ctx)
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return nil, fmt.Errorf("A ScriptPack can't be a dependency")
}

func (a *App) actionRebuild(rctx router.Context) error {
	ctx := rctx.(*app.Context)
	cwd := filepath.Dir(ctx.Appfile.Path)

	// Get the path to the rebuild binary
	path := filepath.Join(ctx.Dir, "rebuild", "rebuild.go")

	// Run it to regenerate the contents
	ctx.Ui.Header("Rebuilding ScriptPack data...")
	cmd := exec.Command("go", "generate")
	cmd.Dir = cwd
	if err := execHelper.Run(ctx.Ui, cmd); err != nil {
		return err
	}

	cmd = exec.Command("go", "run", path)
	cmd.Dir = cwd
	if err := execHelper.Run(ctx.Ui, cmd); err != nil {
		return err
	}

	// Success
	ctx.Ui.Message("ScriptPack data rebuilt!")
	return nil
}

func (a *App) actionTest(rctx router.Context) error {
	ctx := rctx.(*app.Context)

	// Rebuild
	if err := a.actionRebuild(rctx); err != nil {
		return err
	}

	// Verify we have the files
	dir := filepath.Join(filepath.Dir(ctx.Appfile.Path), "_scriptpack_staging")
	if _, err := os.Stat(dir); err != nil {
		return fmt.Errorf(
			"The directory with the built ScriptPack files doesn't exist!\n" +
				"Please build this with `otto dev scriptpack-rebuild` prior to\n" +
				"running tests.")
	}

	// Get the env vars
	f, err := os.Open(filepath.Join(dir, "env"))
	if err != nil {
		return err
	}
	var env map[string]string
	err = json.NewDecoder(f).Decode(&env)
	f.Close()
	if err != nil {
		return err
	}

	// Build the command we execute to run the tests
	cmd := []string{
		"docker",
		"run",
		"-v /vagrant:/devroot",
	}
	for k, v := range env {
		cmd = append(cmd, fmt.Sprintf("-e %s=%s", k, v))
	}
	cmd = append(cmd, "hashicorp/otto-scriptpack-test-ubuntu:14.04")
	cmd = append(cmd, "bats")

	// Determine the test to execute
	testPath := "test"
	if args := rctx.RouteArgs(); len(args) > 0 {
		testPath = args[0]
	}
	if !filepath.IsAbs(testPath) {
		testPath = fmt.Sprintf("/devroot/" + testPath)
	}

	cmd = append(cmd, testPath)

	// Run the command
	ctx.Ui.Header(fmt.Sprintf("Executing: %s", strings.Join(cmd, " ")))
	v := &vagrant.DevOptions{}
	v.Vagrant(ctx).Execute("ssh", "-c", strings.Join(cmd, " "))

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
	actionTestSyn    = "Run ScriptPack tests"
)

const actionRebuildHelp = `
Usage: otto dev scriptpack-rebuild

  Rebuilds the ScriptPack files and dependencies into a single directory
  for dev and test within the development environment.

  This command must be run before running tests after making any changes
  to the ScriptPack.

`

const actionTestHelp = `
Usage: otto dev scriptpack-test [path]

  Tests the ScriptPack with a BATS test file with the given path.
  If no path is specified, the entire "test" directory is used.

  If a path is given, it must be relative to the working directory.
  If an absolute path is given, it must be convertable to a relative
  subpath of this directory.

`
