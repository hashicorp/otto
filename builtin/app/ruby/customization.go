package rubyapp

import (
	"fmt"
	"path/filepath"

	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

const defaultLatestVersion = "2.2"

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) processRuby(d *schema.FieldData) error {
	vsn := d.Get("ruby_version")

	// If we were asked to detect the version, we attempt to do so.
	// If we can't detect it for non-erroneous reasons, we use our default.
	if vsn == "detect" {
		var err error
		c.Opts.Ctx.Ui.Header("Detecting Ruby version to use...")
		vsn, err = detectRubyVersionGemfile(filepath.Dir(c.Opts.Ctx.Appfile.Path))
		if err != nil {
			return err
		}
		if vsn != "" {
			c.Opts.Ctx.Ui.Message(fmt.Sprintf(
				"Detected desired Ruby version: %s", vsn))
		}
		if vsn == "" {
			vsn = defaultLatestVersion
			c.Opts.Ctx.Ui.Message(fmt.Sprintf(
				"No desired Ruby version found! Will use the default: %s", vsn))
		}
	}

	c.Opts.Bindata.Context["ruby_version"] = vsn
	return nil
}
