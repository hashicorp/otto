package goapp

import (
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) processGo(d *schema.FieldData) error {
	c.Opts.Bindata.Context["dev_go_version"] = d.Get("go_version")

	// Go is really finicky about the GOPATH. To help make the dev
	// environment and build environment more correct, we attempt to
	// detect the GOPATH automatically.
	//
	// We use this GOPATH for example in Vagrant to setup the synced
	// folder directly into the GOPATH properly. Magic!
	gopathPath := d.Get("import_path").(string)
	if gopathPath == "" {
		var err error
		c.Opts.Ctx.Ui.Header("Detecting application import path for GOPATH...")
		gopathPath, err = detectImportPath(c.Opts.Ctx)
		if err != nil {
			return err
		}
	}

	folderPath := "/vagrant"
	if gopathPath != "" {
		folderPath = "/opt/gopath/src/" + gopathPath
	}

	c.Opts.Bindata.Context["import_path"] = gopathPath
	c.Opts.Bindata.Context["shared_folder_path"] = folderPath

	return nil
}
