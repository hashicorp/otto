package vagrant

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	execHelper "github.com/hashicorp/otto/helper/exec"
)

// DevOptions is the configuration struct used for Dev.
type DevOptions struct {
	// Dir is the path to the directory with the Vagrantfile. This
	// will default to `#{ctx.Dir}/dev` if empty.
	Dir string

	// Instructions are help text that is shown after creating the
	// development environment.
	Instructions string

	// Deps are the dependencies that we have. This will be used to
	// automatically modify the Vagrantfile.
	Deps []*app.DevDep
}

// Dev can be used as an implementation of app.App.Dev to automatically
// handle creating a development environment and forwarding commands down
// to Vagrant.
func Dev(ctx *app.Context, opts *DevOptions) error {
	switch ctx.Action {
	case "":
		return opts.actionUp(ctx)
	case "destroy":
		return opts.actionDestroy(ctx)
	case "ssh":
		return opts.actionSSH(ctx)
	case "vagrant":
		return opts.actionRaw(ctx)
	default:
		return fmt.Errorf("Unknown action for dev: %s", ctx.Action)
	}
}

func (opts *DevOptions) actionDestroy(ctx *app.Context) error {
	ctx.Ui.Header("Destroying the local development environment...")
	cmd := opts.command(ctx, "destroy", "-f")
	if err := execHelper.Run(ctx.Ui, cmd); err != nil {
		return opts.vagrantError(err)
	}

	ctx.Ui.Raw("\n")
	ctx.Ui.Header("[green]Development environment has been destroyed!")

	return nil
}

func (opts *DevOptions) actionRaw(ctx *app.Context) error {
	ctx.Ui.Header(fmt.Sprintf(
		"Executing: 'vagrant %s'", strings.Join(ctx.ActionArgs, " ")))
	cmd := opts.command(ctx, ctx.ActionArgs...)
	if err := execHelper.Run(ctx.Ui, cmd); err != nil {
		return opts.vagrantError(err)
	}

	return nil
}

func (opts *DevOptions) actionSSH(ctx *app.Context) error {
	ctx.Ui.Header("Executing SSH. This may take a few seconds...")
	cmd := opts.command(ctx, "ssh")
	if err := execHelper.Run(ctx.Ui, cmd); err != nil {
		return opts.vagrantError(err)
	}

	return nil
}

func (opts *DevOptions) actionUp(ctx *app.Context) error {
	// Build the command to execute
	cmd := opts.command(ctx, "up")

	// Output some info the user prior to running
	ctx.Ui.Header(
		"Creating local development environment with Vagrant if it doesn't exist...")
	ctx.Ui.Message(
		"Raw Vagrant output will begin streaming in below. Otto does\n" +
			"not create this output. It is mirrored directly from Vagrant\n" +
			"while the development environment is being created.\n\n")

	// Run it!
	if err := execHelper.Run(ctx.Ui, cmd); err != nil {
		return opts.vagrantError(err)
	}

	// Success, let the user know whats up
	ctx.Ui.Raw("\n")
	ctx.Ui.Header("[green]Development environment successfully created!")
	if opts.Instructions != "" {
		ctx.Ui.Message(opts.Instructions)
	}

	return nil
}

func (opts *DevOptions) command(ctx *app.Context, command ...string) *exec.Cmd {
	// Build the command to execute
	cmd := exec.Command("vagrant", command...)
	cmd.Dir = filepath.Join(ctx.Dir, "dev")
	if opts.Dir != "" {
		cmd.Dir = opts.Dir
	}

	log.Printf("[DEBUG] dev: executing vagrant up in dir: %s", cmd.Dir)
	return cmd
}

func (opts *DevOptions) vagrantError(err error) error {
	return fmt.Errorf(
		"Error executing Vagrant: %s\n\n"+
			"The error messages from Vagrant are usually very informative.\n"+
			"Please read it carefully and fix any issues it mentions. If\n"+
			"the message isn't clear, please report this to the Otto project.",
		err)
}
