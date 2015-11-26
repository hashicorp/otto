package dockerext

import (
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) processDocker(d *schema.FieldData) error {
	image := d.Get("image").(string)
	if image == "" {
		image = c.Opts.Ctx.Application.Name
	}

	c.Opts.Bindata.Context["docker_image"] = image
	c.Opts.Bindata.Context["run_args"] = d.Get("run_args").(string)
	c.Opts.Bindata.Context["command"] = d.Get("command").(string)
	return nil
}
