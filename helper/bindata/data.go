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

	// Template helpers
	_ "github.com/hashicorp/otto/helper/pongo2"
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

		err = d.renderLowLevel(filepath.Join(dst, asset), asset, bytes.NewReader(data))
		if err != nil {
			return err
		}
	}

	return nil
}

// RenderAsset renders a single bindata asset. This file
// will be processed as a template if it ends in ".tpl".
func (d *Data) RenderAsset(dst, src string) error {
	data, err := d.Asset(src)
	if err != nil {
		return err
	}

	return d.renderLowLevel(dst, src, bytes.NewReader(data))
}

// RenderReal renders a real file (not a bindata'd file). This file
// will be processed as a template if it ends in ".tpl".
func (d *Data) RenderReal(dst, src string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	return d.renderLowLevel(dst, src, f)
}

func (d *Data) renderLowLevel(dst string, src string, r io.Reader) error {
	var err error

	// Determine the filename and whether we're dealing with a template
	var tpl *pongo2.Template = nil
	filename := src
	if strings.HasSuffix(filename, ".tpl") {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return err
		}

		filename = strings.TrimSuffix(filename, ".tpl")
		tpl, err = pongo2.FromString(buf.String())
		if err != nil {
			return err
		}
	}

	// Make the directory containing the final path.
	dir := filepath.Dir(dst)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Create the file itself
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	// If it isn't a template, do a direct byte copy
	if tpl == nil {
		_, err = io.Copy(f, r)
		return err
	}

	return tpl.ExecuteWriter(d.Context, f)
}
