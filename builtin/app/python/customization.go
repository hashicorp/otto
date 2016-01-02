package pythonapp

import (
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) process(d *schema.FieldData) error {
	c.Opts.Bindata.Context["python_version"] = d.Get("python_version")
	c.Opts.Bindata.Context["python_entrypoint"] = d.Get("python_entrypoint")
	return nil
}
