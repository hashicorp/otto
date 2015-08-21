package vagrant

import (
	"fmt"
	"log"

	"github.com/hashicorp/otto/app"
)

type BuildOptions struct {
	// Dir is the directory where Vagrant will be executed for the
	// build. This should be the directory with the Vagrantfile, usually.
	Dir string

	// Script is the script to execute within the VM. This script must
	// not ask for input.
	Script string
}

// Build can be used to use Vagrant to build something. This will handle
// starting Vagrant, running the script, collecting files into a list,
// and destroying the Vagrant environment.
func Build(ctx *app.Context, opts *BuildOptions) error {
	log.Printf(
		"[INFO] Vagrant build for '%s' in dir: %s",
		ctx.Appfile.Application.Name, opts.Dir)

	vagrant := &Vagrant{Dir: opts.Dir, Ui: ctx.Ui}

	// tryDestroy is a helper function that we make a local here
	// since we have to clean up in so many potential places.
	tryDestroy := func() error {
		err := vagrant.Execute("destroy", "-f")
		if err != nil {
			ctx.Ui.Header(fmt.Sprintf(
				"[red]Error destroying the Vagrant environment! There may be\n"+
					"lingering resources. The working directory where Vagrant\n"+
					"can be run to check is below. Please manually clean up\n"+
					"the resources:\n\n%s\n\n%s",
				opts.Dir, err))
		}

		return err
	}

	// Bring the environment up. If there is an error, we need to
	// destroy the VM because `vagrant up` can error even after the
	// VM is built.
	if err := vagrant.Execute("up"); err != nil {
		ctx.Ui.Header(fmt.Sprintf(
			"[red]Error while bringing up the Vagrant environment.\n" +
				"The error message will be shown below. First, Otto\n" +
				"will attempt to destroy the machine."))
		tryDestroy()
		return err
	}

	// The environment is running. Execute the build script.
	if err := vagrant.Execute("ssh", "-c", opts.Script); err != nil {
		ctx.Ui.Header(fmt.Sprintf(
			"[red]Error while building in the Vagrant environment!\n" +
				"The error message will be shown below. First, Otto will\n" +
				"attempt to destroy the machine."))
		tryDestroy()
		return err
	}

	// Clean up the Vagrant environment
	ctx.Ui.Message("Vagrant-based build is complete. Deleting Vagrant environment...")
	tryDestroy()
	return nil
}
