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
	c.Opts.Bindata.Context["scala_version"] = d.Get("scala_version")
	c.Opts.Bindata.Context["sbt_version"] = d.Get("sbt_version")
	c.Opts.Bindata.Context["lein_version"] = d.Get("lein_version")
	return nil
}
