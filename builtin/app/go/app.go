package goapp

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/helper/bindata"
	execHelper "github.com/hashicorp/otto/helper/exec"
)

//go:generate go-bindata -pkg=goapp -nomemcopy ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	data := &bindata.Data{
		Asset:    Asset,
		AssetDir: AssetDir,
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

	return nil, nil
}

func (a *App) Dev(ctx *app.Context) error {
	// Build the command to execute
	cmd := exec.Command("vagrant", "up")
	cmd.Dir = filepath.Join(ctx.Dir, "dev")

	// Output some info the user prior to running
	ctx.Ui.Header("Executing Vagrant to manage local dev environment...")
	ctx.Ui.Message(
		"Raw Vagrant output will begin streaming in below. Otto does\n" +
			"not create this output. It is mirrored directly from Vagrant\n" +
			"while the development environment is being created.\n\n")

	// Run it!
	if err := execHelper.Run(ctx.Ui, cmd); err != nil {
		return fmt.Errorf(
			"Error executing Vagrant: %s\n\n" +
				"The error messages from Vagrant are usually very informative.\n" +
				"Please read it carefully and fix any issues it mentions. If\n" +
				"the message isn't clear, please report this to the Otto project.")
	}

	// Success, let the user know whats up
	ctx.Ui.Header("[green]Development environment successfully created!")

	return nil
}
