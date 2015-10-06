package bindata

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2"

	// Template helpers
	_ "github.com/hashicorp/otto/helper/pongo2"
)

//go:generate go-bindata -o=bindata_test.go -pkg=bindata -nomemcopy -nometadata ./test-data/...

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

	// SharedExtends is a mapping of share prefixes and files that can be
	// accessed using {% extends %} in templates. Example:
	// {% extends "foo:bar/baz.tpl" %} would find the "bar/baz.tpl" in the
	// "foo" share.
	SharedExtends map[string]*Data
}

// CopyDir copies all the assets from the given prefix to the destination
// directory. It will automatically set file permissions, create folders,
// etc.
func (d *Data) CopyDir(dst, prefix string) error {
	return d.copyDir(dst, prefix)
}

// RenderAsset renders a single bindata asset. This file
// will be processed as a template if it ends in ".tpl".
func (d *Data) RenderAsset(dst, src string) error {
	data, err := d.Asset(src)
	if err != nil {
		return err
	}

	return d.renderLowLevel(dst, src, "", bytes.NewReader(data))
}

// RenderReal renders a real file (not a bindata'd file). This file
// will be processed as a template if it ends in ".tpl".
func (d *Data) RenderReal(dst, src string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	return d.renderLowLevel(dst, src, "", f)
}

// RenderString renders a string.
func (d *Data) RenderString(tpl string) (string, error) {
	// Make a temporary file for the contents. This is kind of silly we
	// need to do this but we can make this render in-memory later.
	tf, err := ioutil.TempFile("", "otto")
	if err != nil {
		return "", err
	}
	tf.Close()
	defer os.Remove(tf.Name())

	// Render
	err = d.renderLowLevel(tf.Name(), "dummy.tpl", "", strings.NewReader(tpl))
	if err != nil {
		return "", err
	}

	// Copy the file contents back into memory
	var result bytes.Buffer
	f, err := os.Open(tf.Name())
	if err != nil {
		return "", err
	}
	_, err = io.Copy(&result, f)
	f.Close()
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func (d *Data) copyDir(dst, prefix string) error {
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

			if err := d.copyDir(filepath.Join(dst, asset), assetFull); err != nil {
				return err
			}

			continue
		}

		err = d.renderLowLevel(filepath.Join(dst, asset), asset, prefix, bytes.NewReader(data))
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *Data) renderLowLevel(dst string, src string, prefix string, r io.Reader) error {
	var err error

	// Determine the filename and whether we're dealing with a template
	var tpl *pongo2.Template = nil
	filename := src
	if strings.HasSuffix(filename, ".tpl") {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err != nil {
			return err
		}

		base := filepath.Dir(filename)
		if prefix != "" {
			base = filepath.Join(prefix, base)
		}

		// Create the template set so we can control loading
		tplSet := pongo2.NewSet("otto", &tplLoader{
			Data: d,
			Base: base,
		})

		// Parse the template
		dst = strings.TrimSuffix(dst, ".tpl")
		tpl, err = tplSet.FromString(buf.String())
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
