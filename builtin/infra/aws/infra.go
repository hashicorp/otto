package aws

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/infrastructure"
)

//go:generate go-bindata -pkg=aws -nomemcopy ./data/...

// Infra is an implementation of infrastructure.Infrastructure
type Infra struct{}

func (i *Infra) Execute(ctx *infrastructure.Context) error {
	statePath := filepath.Join(ctx.Dir, "terraform.tfstate.new")

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

	// Start the Terraform command
	if err := cmd.Run(); err != nil {
		return fmt.Errorf(
			"Error running Terraform: %s", err)
	}

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
