package custom

import (
	"path/filepath"

	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) processBuild(d *schema.FieldData) error {
	p, ok := d.GetOk("packer")
	if !ok {
		return nil
	}

	c.Opts.Bindata.Context["build_packer_path"] = p.(string)
	c.Opts.Callbacks = append(c.Opts.Callbacks, c.compileCustomBuild(d))
	return nil
}

func (c *customizations) processDeploy(d *schema.FieldData) error {
	tf, ok := d.GetOk("terraform")
	if !ok {
		return nil
	}

	c.Opts.Bindata.Context["deploy_terraform_path"] = tf.(string)
	c.Opts.Callbacks = append(c.Opts.Callbacks, c.compileCustomDeploy(d))
	return nil
}

func (c *customizations) processDev(d *schema.FieldData) error {
	p, ok := d.GetOk("vagrant")
	if !ok {
		return nil
	}

	c.Opts.Bindata.Context["dev_vagrant_path"] = p.(string)
	c.Opts.Callbacks = append(c.Opts.Callbacks, c.compileCustomDev(d))
	return nil
}

func (c *customizations) processDevDep(d *schema.FieldData) error {
	if _, ok := d.GetOk("vagrantfile"); !ok {
		return nil
	}

	c.Opts.Callbacks = append(c.Opts.Callbacks, c.compileCustomDevDep(d))
	return nil
}

func (c *customizations) compileCustomBuild(d *schema.FieldData) compile.CompileCallback {
	return func() error {
		return c.Opts.Bindata.RenderAsset(
			filepath.Join(c.Opts.Ctx.Dir, "build", "packer_path"),
			"data/sentinels/packer_path.tpl")
	}
}

func (c *customizations) compileCustomDeploy(d *schema.FieldData) compile.CompileCallback {
	return func() error {
		return c.Opts.Bindata.RenderAsset(
			filepath.Join(c.Opts.Ctx.Dir, "deploy", "terraform_path"),
			"data/sentinels/terraform_path.tpl")
	}
}

func (c *customizations) compileCustomDev(d *schema.FieldData) compile.CompileCallback {
	return func() error {
		return c.Opts.Bindata.RenderAsset(
			filepath.Join(c.Opts.Ctx.Dir, "dev", "vagrantfile_path"),
			"data/sentinels/vagrant_path.tpl")
	}
}

func (c *customizations) compileCustomDevDep(d *schema.FieldData) compile.CompileCallback {
	vf := d.Get("vagrantfile").(string)

	return func() error {
		data := c.Opts.Bindata
		fragment := data.Context["fragment_path"].(string)
		if err := data.RenderReal(fragment, vf); err != nil {
			return err
		}

		return data.RenderAsset(
			filepath.Join(c.Opts.Ctx.Dir, "dev", "Vagrantfile"),
			"data/dev/Vagrantfile.tpl")
	}
}
