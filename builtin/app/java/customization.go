package javaapp

import (
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) processDev(d *schema.FieldData) error {
	c.Opts.Bindata.Context["gradle_version"] = d.Get("gradle_version")
	c.Opts.Bindata.Context["maven_version"] = d.Get("maven_version")
	return nil
}
