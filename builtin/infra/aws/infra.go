package aws

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/otto/infrastructure"
	"github.com/mitchellh/iochan"
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
	cmd.Stdout = out_w
	cmd.Stderr = out_w

	ctx.Ui.Header("Executing Terraform...")
	ctx.Ui.Message(
		"Raw Terraform output will begin streaming in below. Otto\n" +
			"does not create this output. It is mirrored directly from\n" +
			"Terraform while the infrastructure is being created.\n\n")

	// Start the Terraform command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf(
			"Error starting Terraform: %s", err)
	}

	// Execute and copy output
	for line := range iochan.DelimReader(out_r, '\n') {
		ctx.Ui.Raw(line)
	}

	return nil
}

func (i *Infra) Compile(ctx *infrastructure.Context) (*infrastructure.CompileResult, error) {
	// Create the output directory
	if err := os.MkdirAll(ctx.Dir, 0755); err != nil {
		return nil, err
	}

	// Get all the assets in our flavor directory and process them all
	// into the output directory.
	prefix := "data/" + ctx.Infra.Flavor
	assets, err := AssetDir(prefix)
	if err != nil {
		return nil, err
	}

	for _, asset := range assets {
		log.Printf("[DEBUG] Writing file: %s", asset)

		data := MustAsset(prefix + "/" + asset)

		// If we have a parent directory create that
		if dir := filepath.Dir(asset); dir != "." {
			// TODO
			panic("TODO")
		}

		// Write the file itself
		f, err := os.Create(filepath.Join(ctx.Dir, asset))
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(f, bytes.NewReader(data))
		f.Close()
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *Infra) Flavors() []string {
	return nil
}
