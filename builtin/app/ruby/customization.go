package rubyapp

import (
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) processRuby(d *schema.FieldData) error {
	c.Opts.Bindata.Context["ruby_version"] = d.Get("ruby_version")
	return nil
}
