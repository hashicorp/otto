package vagrant

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/helper/router"
)

// DevOptions is the configuration struct used for Dev.
type DevOptions struct {
	// Dir is the path to the directory with the Vagrantfile. This
	// will default to `#{ctx.Dir}/dev` if empty.
	Dir string

	// DataDir is the path to the directory where Vagrant should store its data.
	// Defaults to `#{ctx.LocalDir/vagrant}` if empty.
	DataDir string

	// Layer, if non-nil, will be the set of layers that this environment
	// builds on top of. If this is set, then the layers will be managed
	// automatically by this.
	//
	// If this is nil, then layers won't be used.
	Layer *Layered

	// Instructions are help text that is shown after creating the
	// development environment.
	Instructions string
}

// Dev can be used as an implementation of app.App.Dev to automatically
// handle creating a development environment and forwarding commands down
// to Vagrant.
func Dev(opts *DevOptions) *router.Router {
	return &router.Router{
		Actions: map[string]router.Action{
			"": &router.SimpleAction{
				ExecuteFunc:  opts.actionUp,
				SynopsisText: actionUpSyn,
				HelpText:     strings.TrimSpace(actionUpHelp),
			},

			"address": &router.SimpleAction{
				ExecuteFunc:  opts.actionAddress,
				SynopsisText: actionAddressSyn,
				HelpText:     strings.TrimSpace(actionAddressHelp),
			},

			"destroy": &router.SimpleAction{
				ExecuteFunc:  opts.actionDestroy,
				SynopsisText: actionDestroySyn,
				HelpText:     strings.TrimSpace(actionDestroyHelp),
			},

			"halt": &router.SimpleAction{
				ExecuteFunc:  opts.actionHalt,
				SynopsisText: actionHaltSyn,
				HelpText:     strings.TrimSpace(actionHaltHelp),
			},

			"layers": &router.SimpleAction{
				ExecuteFunc:  opts.actionLayers,
				SynopsisText: actionLayersSyn,
				HelpText:     strings.TrimSpace(actionLayersHelp),
			},

			"ssh": &router.SimpleAction{
				ExecuteFunc:  opts.actionSSH,
				SynopsisText: actionSSHSyn,
				HelpText:     strings.TrimSpace(actionSSHHelp),
			},

			"vagrant": &router.SimpleAction{
				ExecuteFunc:  opts.actionRaw,
				SynopsisText: actionVagrantSyn,
				HelpText:     strings.TrimSpace(actionVagrantHelp),
			},
		},
	}
}

func (opts *DevOptions) actionAddress(rctx router.Context) error {
	ctx := rctx.(*app.Context)
	ctx.Ui.Raw(ctx.DevIPAddress + "\n")
	return nil
}

func (opts *DevOptions) actionDestroy(rctx router.Context) error {
	ctx := rctx.(*app.Context)
	project := Project(&ctx.Shared)
	if err := project.InstallIfNeeded(); err != nil {
		return err
	}

	ctx.Ui.Header("Destroying the local development environment...")
	vagrant := opts.vagrant(ctx)

	// If the Vagrant directory doesn't exist, then we're already deleted.
	// So we just verify here that it exists and then call destroy only
	// if it does.
	log.Printf("[DEBUG] vagrant: verifying data dir exists: %s", vagrant.DataDir)
	_, err := os.Stat(vagrant.DataDir)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("[ERROR] vagrant: err: %s", err)
		return err
	}
	if err == nil {
		if err := vagrant.Execute("destroy", "-f"); err != nil {
			return err
		}
		ctx.Ui.Raw("\n")
	}

	// Store the dev status into the directory. We just do this before
	// since there are a lot of cases where Vagrant fails but still imported.
	// We just override any prior dev.
	ctx.Ui.Header("Deleting development environment metadata...")
	if opts.Layer != nil {
		if err := opts.Layer.RemoveEnv(vagrant); err != nil {
			return fmt.Errorf(
				"Error preparing dev environment: %s", err)
		}
	}

	if err := ctx.Directory.DeleteDev(opts.devLookup(ctx)); err != nil {
		return fmt.Errorf(
			"Error deleting dev environment metadata: %s", err)
	}

	if err := opts.sshCache(ctx).Delete(); err != nil {
		return fmt.Errorf(
			"Error cleaning SSH cache: %s", err)
	}

	ctx.Ui.Header("[green]Development environment has been destroyed!")
	return nil
}

func (opts *DevOptions) actionHalt(rctx router.Context) error {
	ctx := rctx.(*app.Context)
	project := Project(&ctx.Shared)
	if err := project.InstallIfNeeded(); err != nil {
		return err
	}

	ctx.Ui.Header("Halting the the local development environment...")

	if err := opts.vagrant(ctx).Execute("halt"); err != nil {
		return err
	}

	ctx.Ui.Header("[green]Development environment halted!")

	return nil
}

func (opts *DevOptions) actionRaw(rctx router.Context) error {
	ctx := rctx.(*app.Context)
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

func (opts *DevOptions) actionSSH(rctx router.Context) error {
	ctx := rctx.(*app.Context)

	dev, err := ctx.Directory.GetDev(opts.devLookup(ctx))
	if err != nil {
		return err
	}
	if dev == nil {
		return fmt.Errorf(
			"The development environment hasn't been created yet! Please\n" +
				"create the development environmet by running `otto dev` before\n" +
				"attempting to SSH.")
	}

	project := Project(&ctx.Shared)
	if err := project.InstallIfNeeded(); err != nil {
		return err
	}

	ctx.Ui.Header("Executing SSH. This may take a few seconds...")
	return opts.sshCache(ctx).Exec(true)
}

func (opts *DevOptions) actionUp(rctx router.Context) error {
	ctx := rctx.(*app.Context)
	project := Project(&ctx.Shared)
	if err := project.InstallIfNeeded(); err != nil {
		return err
	}

	// If we are layered, then let the user know we're going to use
	// a layer development environment...
	if opts.Layer != nil {
		pending, err := opts.Layer.Pending()
		if err != nil {
			return fmt.Errorf("Error checking dev layer status: %s", err)
		}

		if len(pending) > 0 {
			ctx.Ui.Header("Creating development environment layers...")
			ctx.Ui.Message(
				"Otto uses layers to speed up building development environments.\n" +
					"Each layer only needs to be built once. We've detected that the\n" +
					"layers below aren't created yet. These will be built this time.\n" +
					"Future development envirionments will use the cached versions\n" +
					"to be much, much faster.")
		}

		if err := opts.Layer.Build(&ctx.Shared); err != nil {
			return fmt.Errorf(
				"Error building dev environment layers: %s", err)
		}
	}

	// Output some info the user prior to running
	ctx.Ui.Header(
		"Creating local development environment with Vagrant if it doesn't exist...")

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
	vagrant := opts.vagrant(ctx)
	if opts.Layer != nil {
		if err := opts.Layer.ConfigureEnv(vagrant); err != nil {
			return fmt.Errorf(
				"Error preparing dev environment: %s", err)
		}

		// Configure the environment as ready
		if err := opts.Layer.SetEnv(vagrant, envStateReady); err != nil {
			return fmt.Errorf(
				"Error preparing dev environment: %s", err)
		}
	}
	if err := vagrant.Execute("up"); err != nil {
		return err
	}

	// Cache the SSH info
	ctx.Ui.Header("Caching SSH credentials from Vagrant...")
	if err := opts.sshCache(ctx).Cache(); err != nil {
		return err
	}

	// Success, let the user know whats up
	ctx.Ui.Header("[green]Development environment successfully created!")
	ctx.Ui.Message(fmt.Sprintf("IP address: %s", ctx.DevIPAddress))
	if opts.Instructions != "" {
		ctx.Ui.Message("\n" + opts.Instructions)
	}

	return nil
}

func (opts *DevOptions) actionLayers(rctx router.Context) error {
	if opts.Layer == nil {
		return fmt.Errorf(
			"This development environment does not use layers.\n" +
				"This command can only be used to manage development\n" +
				"environments with layers.")
	}

	ctx := rctx.(*app.Context)
	fs := flag.NewFlagSet("otto", flag.ContinueOnError)
	graph := fs.Bool("graph", false, "show graph")
	prune := fs.Bool("prune", false, "prune unused layers")
	if err := fs.Parse(rctx.RouteArgs()); err != nil {
		return err
	}

	// Graph?
	if *graph {
		graph, err := opts.Layer.Graph()
		if err != nil {
			return err
		}

		ctx.Ui.Raw(graph.String() + "\n")
		return nil
	}

	// Prune?
	if *prune {
		ctx.Ui.Header("Pruning any outdated or unused layers...")
		count, err := opts.Layer.Prune(&ctx.Shared)
		if err != nil {
			return err
		}
		if count == 0 {
			ctx.Ui.Message("No outdated or unused layers were found!")
		} else {
			ctx.Ui.Message(fmt.Sprintf(
				"[green]Pruned %d outdated or unused layers!", count))
		}

		return nil
	}

	// We're just listing the layers. Eventually we probably should
	// output status or something more useful here.
	for _, l := range opts.Layer.Layers {
		ctx.Ui.Raw(l.ID + "\n")
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

func (opts *DevOptions) devLookup(ctx *app.Context) *directory.Dev {
	return &directory.Dev{Lookup: directory.Lookup{AppID: ctx.Appfile.ID}}
}

func (opts *DevOptions) sshCache(ctx *app.Context) *SSHCache {
	return &SSHCache{
		Path:    filepath.Join(ctx.CacheDir, "dev_ssh_cache"),
		Vagrant: opts.vagrant(ctx),
	}
}

// Synopsis text for actions
const (
	actionAddressSyn = "Shows the address to reach the development environment"
	actionUpSyn      = "Starts the development environment"
	actionDestroySyn = "Destroy the development environment"
	actionHaltSyn    = "Halts the development environment"
	actionLayersSyn  = "Manage the layers of this development environment"
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

const actionHaltHelp = `
Usage: otto dev halt

  Halts the development environment.

  This command will stop the development environment. The environment can then
  be started again with 'otto dev'.

`

const actionLayersHelp = `
Usage: otto dev layers [options]

  Manage the development environment layers.

  WARNING: This is an advanced, low level command. You shouldn't need this
  command. It is meant to give you the ability to get out of a bad situation
  if Otto mis-manages your layers. If you run into a scenario where you need
  to use this, please report a bug to Otto so we can think of others ways
  around it.

  This command will manage the layers of the development environment.
  Otto uses layers as a mechanism for caching parts of the development
  environment that aren't often updated. This makes "otto dev" faster
  after the first call.

  If no options are given, the layers will be listed that this development
  environment uses. If multiple conflicting options are given, the first
  in alphabetical order is processed. For example, if both "-graph" and
  "-prune" are specified, the graph will be shown.

Options:

  -graph       Show the full layer graph for Otto
  -prune       Delete all unused or outdated layers

`

const actionSSHHelp = `
Usage: otto dev ssh

  Connect to the development environment via SSH.

  The development environment typically is headless, meaning that the
  preferred way to access it is SSH. This command will automatically SSH
  you into the development environment.

`

const actionAddressHelp = `
Usage: otto dev address

  Output the address to connect to the development environment.

  The development environment is configured with a static IP address.
  This command outputs that address so you can reach it. If you want to
  SSH into the development environment, use 'otto dev ssh'. This address
  is meant for reaching running services such as in a web browser.

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
