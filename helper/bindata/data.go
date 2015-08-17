package bindata

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/flosch/pongo2.v3"
)

//go:generate go-bindata -o=bindata_test.go -pkg=bindata -nomemcopy ./test-data/...

// Data is the struct that wraps the assets go-bindata generates in your
// package to provide more helpers.
type Data struct {
	// Asset, AssetDir are functions used to look up the assets.
	// These match the function signatures used by go-bindata so you
	// can just use method handles for these.
	Asset    func(string) ([]byte, error)
	AssetDir func(string) ([]string, error)

	// Context is the template context that is given when rendering
	Context map[string]interface{}
}

// CopyDir copies all the assets from the given prefix to the destination
// directory. It will automatically set file permissions, create folders,
// etc.
func (d *Data) CopyDir(dst, prefix string) error {
	log.Printf("[DEBUG] Copying all assets: %s => %s", prefix, dst)

	// Get all the assets in the directory
	assets, err := d.AssetDir(prefix)
	if err != nil {
		return err
	}

	// If the destination directory doesn't exist, make that
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// Go through each asset, and copy it into place
	for _, asset := range assets {
		log.Printf("[DEBUG] Copying asset: %s", asset)

		// Load the asset bytes
		assetFull := prefix + "/" + asset
		data, err := d.Asset(assetFull)
		if err != nil {
			// Asset not found... check if it is a directory. If it is
			// a directory, we recursively call CopyDir.
			if _, err := d.AssetDir(assetFull); err != nil {
				return fmt.Errorf("error loading asset %s: %s", asset, err)
			}

			if err := d.CopyDir(filepath.Join(dst, asset), assetFull); err != nil {
				return err
			}

			continue
		}

		// Determine the filename and whether we're dealing with a template
		var tpl *pongo2.Template = nil
		filename := asset
		if strings.HasSuffix(filename, ".tpl") {
			filename = strings.TrimSuffix(filename, ".tpl")
			tpl, err = pongo2.FromString(string(data))
			if err != nil {
				return err
			}
		}

		// Write the file
		f, err := os.Create(filepath.Join(dst, filename))
		if err != nil {
			return err
		}

		// If it isn't a template, do a direct byte copy
		if tpl == nil {
			_, err = io.Copy(f, bytes.NewReader(data))
			f.Close()
			if err != nil {
				return err
			}

			continue
		}

		err = tpl.ExecuteWriter(d.Context, f)
		f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
