package phpapp

import (
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) processPhp(d *schema.FieldData) error {
	c.Opts.Bindata.Context["php_version"] = d.Get("php_version")
	return nil
}
