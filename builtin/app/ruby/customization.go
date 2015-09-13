package rubyapp

import (
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) processDev(d *schema.FieldData) error {
	c.Opts.Bindata.Context["dev_ruby_version"] = d.Get("ruby_version")
	return nil
}
