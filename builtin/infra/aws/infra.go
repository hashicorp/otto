package aws

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/infrastructure"
)

//go:generate go-bindata -pkg=aws -nomemcopy ./data/...

// Infra is an implementation of infrastructure.Infrastructure
type Infra struct{}

func (i *Infra) Execute(ctx *infrastructure.Context) error {
	statePath, err := filepath.Abs(filepath.Join(ctx.Dir, "terraform.tfstate.new"))
	if err != nil {
		return fmt.Errorf(
			"Error building state output path: %s\n\n"+
				"This is an internal error that should really never happen.\n"+
				"No infrastructure was created. Please report this as a bug.", err)
	}

	// Build the command to execute
	out_r, out_w := io.Pipe()
	cmd := exec.Command(
		"terraform",
		"apply",
		"-state-out", statePath)
	cmd.Dir = ctx.Dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = out_w
	cmd.Stderr = out_w

	ctx.Ui.Header("Executing Terraform to build infrastructure...")
	ctx.Ui.Message(
		"Raw Terraform output will begin streaming in below. Otto\n" +
			"does not create this output. It is mirrored directly from\n" +
			"Terraform while the infrastructure is being created.\n\n" +
			"Terraform may ask for input. For infrastructure provider\n" +
			"credentials, be sure to enter the same credentials\n" +
			"consistently within the same Otto environment." +
			"\n\n")

	// Copy output to the UI until we can't
	go func() {
		defer out_w.Close()
		var buf [1024]byte
		for {
			n, err := out_r.Read(buf[:])
			if n > 0 {
				ctx.Ui.Raw(string(buf[:n]))
			}

			// We just break on any error. io.EOF is not an error and
			// is our true exit case, but any other error we don't really
			// handle here. It probably means something went wrong
			// somewhere else anyways.
			if err != nil {
				break
			}
		}
	}()

	var infra directory.Infra
	infra.State = directory.InfraStateReady

	// Start the Terraform command
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("Error running Terraform: %s", err)
		infra.State = directory.InfraStatePartial
	}

	// Read the outputs if everything is looking good so far
	if err == nil {
		infra.Outputs, err = terraform.Outputs(statePath)
		if err != nil {
			err = fmt.Errorf("Error reading Terraform outputs: %s", err)
			infra.State = directory.InfraStatePartial
		}
	}

	// Save the infrastructure information
	ctx.Ui.Header("Terraform execution complete. Saving results...")
	if err := ctx.Directory.PutInfra("TODO: ID", &infra); err != nil {
		return fmt.Errorf(
			"Error storing infrastructure data: %s\n\n"+
				"This means that Otto won't be able to know that your infrastructure\n"+
				"was successfully created. Otto tries a few times to save the\n"+
				"infrastructure. At this point in time, Otto doesn't support gracefully\n"+
				"recovering from this error. Your infrastructure is now orphaned from\n"+
				"Otto's management. Please reference the community for help.\n\n"+
				"A future version of Otto will resolve this.",
			err)
	}

	// If there was an error during the process, then return that.
	if err != nil {
		return fmt.Errorf("Error reading Terraform outputs: %s\n\n"+
			"In this case, Otto is unable to consider the infrastructure ready.\n"+
			"Otto won't lose your infrastructure information. You may just need\n"+
			"to run `otto infra` again and it may work. If this problem persists,\n"+
			"please see the error message and consult the community for help.",
			err)
	}

	// Output something to the user so they know what is going on.
	ctx.Ui.Header("[green]Infrastructure successfully created!")
	ctx.Ui.Message(
		"[green]The infrastructure necessary to deploy this application\n" +
			"is now available. You can now deploy using `otto deploy`.")

	return nil
}

func (i *Infra) Compile(ctx *infrastructure.Context) (*infrastructure.CompileResult, error) {
	data := &bindata.Data{
		Asset:    Asset,
		AssetDir: AssetDir,
	}

	if err := data.CopyDir(ctx.Dir, "data/"+ctx.Infra.Flavor); err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Infra) Flavors() []string {
	return nil
}
