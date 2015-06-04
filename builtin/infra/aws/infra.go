package aws

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hashicorp/otto/infrastructure"
)

//go:generate go-bindata -pkg=aws -nomemcopy ./data/...

// Infra is an implementation of infrastructure.Infrastructure
type Infra struct{}

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
