package vagrant

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/directory"
)

// DevOptions is the configuration struct used for Dev.
type DevOptions struct {
	// Dir is the path to the directory with the Vagrantfile. This
	// will default to `#{ctx.Dir}/dev` if empty.
	Dir string

	// DataDir is the path to the directory where Vagrant should store its data.
	// Defaults to `#{ctx.LocalDir/vagrant}` if empty.
	DataDir string

	// Instructions are help text that is shown after creating the
	// development environment.
	Instructions string
}

// Dev can be used as an implementation of app.App.Dev to automatically
// handle creating a development environment and forwarding commands down
// to Vagrant.
func Dev(opts *DevOptions) *app.Router {
	return &app.Router{
		Actions: map[string]*app.Action{
			"": &app.Action{
				Execute:  opts.actionUp,
				Synopsis: actionUpSyn,
				Help:     strings.TrimSpace(actionUpHelp),
			},

			"destroy": &app.Action{
				Execute:  opts.actionDestroy,
				Synopsis: actionDestroySyn,
				Help:     strings.TrimSpace(actionDestroyHelp),
			},

			"ssh": &app.Action{
				Execute:  opts.actionSSH,
				Synopsis: actionSSHSyn,
				Help:     strings.TrimSpace(actionSSHHelp),
			},

			"vagrant": &app.Action{
				Execute:  opts.actionRaw,
				Synopsis: actionVagrantSyn,
				Help:     strings.TrimSpace(actionVagrantHelp),
			},
		},
	}
}

func (opts *DevOptions) actionDestroy(ctx *app.Context) error {
	project := Project(&ctx.Shared)
	if err := project.InstallIfNeeded(); err != nil {
		return err
	}

	ctx.Ui.Header("Destroying the local development environment...")

	if err := opts.vagrant(ctx).Execute("destroy", "-f"); err != nil {
		return err
	}
	ctx.Ui.Raw("\n")

	// Store the dev status into the directory. We just do this before
	// since there are a lot of cases where Vagrant fails but still imported.
	// We just override any prior dev.
	ctx.Ui.Header("Deleting development environment metadata...")
	dev := &directory.Dev{Lookup: directory.Lookup{AppID: ctx.Appfile.ID}}
	if err := ctx.Directory.DeleteDev(dev); err != nil {
		return fmt.Errorf(
			"Error deleting dev environment metadata: %s", err)
	}

	ctx.Ui.Header("[green]Development environment has been destroyed!")
	return nil
}

func (opts *DevOptions) actionRaw(ctx *app.Context) error {
	project := Project(&ctx.Shared)
	if err := project.InstallIfNeeded(); err != nil {
		return err
	}

	ctx.Ui.Header(fmt.Sprintf(
		"Executing: 'vagrant %s'", strings.Join(ctx.ActionArgs, " ")))

	if err := opts.vagrant(ctx).Execute(ctx.ActionArgs...); err != nil {
		return err
	}

	return nil
}

func (opts *DevOptions) actionSSH(ctx *app.Context) error {
	project := Project(&ctx.Shared)
	if err := project.InstallIfNeeded(); err != nil {
		return err
	}

	ctx.Ui.Header("Executing SSH. This may take a few seconds...")
	if err := opts.vagrant(ctx).Execute("ssh"); err != nil {
		return err
	}

	return nil
}

func (opts *DevOptions) actionUp(ctx *app.Context) error {
	project := Project(&ctx.Shared)
	if err := project.InstallIfNeeded(); err != nil {
		return err
	}

	// Output some info the user prior to running
	ctx.Ui.Header(
		"Creating local development environment with Vagrant if it doesn't exist...")
	ctx.Ui.Message(
		"Raw Vagrant output will begin streaming in below. Otto does\n" +
			"not create this output. It is mirrored directly from Vagrant\n" +
			"while the development environment is being created.\n\n")

	// Store the dev status into the directory. We just do this before
	// since there are a lot of cases where Vagrant fails but still imported.
	// We just override any prior dev.
	dev := &directory.Dev{Lookup: directory.Lookup{AppID: ctx.Appfile.ID}}
	dev.MarkReady()
	if err := ctx.Directory.PutDev(dev); err != nil {
		return fmt.Errorf(
			"Error saving dev environment metadata: %s", err)
	}

	// Run it!
	if err := opts.vagrant(ctx).Execute("up"); err != nil {
		return err
	}
	// Success, let the user know whats up
	ctx.Ui.Raw("\n")
	ctx.Ui.Header("[green]Development environment successfully created!")
	if opts.Instructions != "" {
		ctx.Ui.Message(opts.Instructions)
	}

	return nil
}

func (opts *DevOptions) vagrant(ctx *app.Context) *Vagrant {
	dir := opts.Dir
	if dir == "" {
		dir = filepath.Join(ctx.Dir, "dev")
	}
	dataDir := opts.DataDir
	if dataDir == "" {
		dataDir = filepath.Join(ctx.LocalDir, "vagrant")
	}
	return &Vagrant{
		Dir:     dir,
		DataDir: dataDir,
		Ui:      ctx.Ui,
	}
}

// Synopsis text for actions
const (
	actionUpSyn      = "Starts the development environment"
	actionDestroySyn = "Destroy the development environment"
	actionSSHSyn     = "SSH into the development environment"
	actionVagrantSyn = "Run arbitrary Vagrant commands"
)

// Help text for actions
const actionUpHelp = `
Usage: otto dev

  Builds and starts the development environment.

  The development environment runs locally via Vagrant. Otto manages
  Vagrant for you. All upstream dependencies will automatically be started
  and running within the development environment.

  At the end of running this command, help text will be shown that tell
  you how to interact with the build environment.
`

const actionDestroyHelp = `
Usage: otto dev destroy

  Destroys the development environment.

  This command will stop and delete the development environment.
  Any data that was put onto the development environment will be deleted,
  except for your own project's code (the directory and any subdirectories
  where the Appfile exists).

`

const actionSSHHelp = `
Usage: otto dev ssh

  Connect to the development environment via SSH.

  The development environment typically is headless, meaning that the
  preferred way to access it is SSH. This command will automatically SSH
  you into the development environment.

`

const actionVagrantHelp = `
Usage: otto dev vagrant [command...]

  Run arbitrary Vagrant commands against the development environment.

  This is for advanced users who know and are comfortable with Vagrant.
  In average day to day usage, this command isn't needed.

  Because the development environment is backed by Vagrant, this command
  lets you access it directly via Vagrant. For example, if you want to
  run "vagrant ssh-config" against the environment, you can use
  "otto dev vagrant ssh-config"

`
