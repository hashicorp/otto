package custom

import (
	"path/filepath"

	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
)

type customizations struct {
	Opts *compile.AppOptions
}

func (c *customizations) process(d *schema.FieldData) error {
	if p, ok := d.GetOk("packer"); ok {
		c.Opts.Bindata.Context["build_packer_path"] = p.(string)
		c.Opts.Callbacks = append(c.Opts.Callbacks, c.compileCustomBuild(d))
	}

	if tf, ok := d.GetOk("terraform"); ok {
		c.Opts.Bindata.Context["deploy_terraform_path"] = tf.(string)
		c.Opts.Callbacks = append(c.Opts.Callbacks, c.compileCustomDeploy(d))
	}

	if p, ok := d.GetOk("dev_vagrantfile"); ok {
		c.Opts.Bindata.Context["dev_vagrant_path"] = p.(string)
		c.Opts.Callbacks = append(c.Opts.Callbacks, c.compileCustomDev(d))
	}

	if _, ok := d.GetOk("dep_vagrantfile"); ok {
		c.Opts.Callbacks = append(c.Opts.Callbacks, c.compileCustomDevDep(d))
	}

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
		c.Opts.Bindata.RenderReal(
			filepath.Join(c.Opts.Ctx.Dir, "dev", "Vagrantfile"),
			c.Opts.Bindata.Context["dev_vagrant_path"].(string))
		return c.Opts.Bindata.RenderAsset(
			filepath.Join(c.Opts.Ctx.Dir, "dev", "vagrantfile_path"),
			"data/sentinels/vagrant_path.tpl")
	}
}

func (c *customizations) compileCustomDevDep(d *schema.FieldData) compile.CompileCallback {
	vf := d.Get("vagrantfile").(string)

	return func() error {
		if !filepath.IsAbs(vf) {
			vf = filepath.Join(filepath.Dir(c.Opts.Ctx.Appfile.Path), vf)
		}

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
