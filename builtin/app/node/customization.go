package nodeapp

import (
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
	"fmt"
	"path/filepath"	
)

const defaultLatestVersion = "4.1.0"

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) process(d *schema.FieldData) error {

	vsn := d.Get("node_version")

	// If we were asked to detect the version, we attempt to do so.
	// If we can't detect it for non-erroneous reasons, we use our default.
	if vsn == "detect" {
		var err error
		c.Opts.Ctx.Ui.Header("Detecting Node version to use...")
		vsn, err = detectNodeVersionJsonfile(filepath.Dir(c.Opts.Ctx.Appfile.Path))
		if err != nil {
			return err
		}
		if vsn != "" {
			c.Opts.Ctx.Ui.Message(fmt.Sprintf(
				"Detected desired Node version: %s", vsn))
		}
		if vsn == "" {
			vsn = defaultLatestVersion
			c.Opts.Ctx.Ui.Message(fmt.Sprintf(
				"No desired Node version found! Will use the default: %s", vsn))
		}
	}

	c.Opts.Bindata.Context["node_version"] = vsn
	return nil
}
