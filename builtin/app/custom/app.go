package custom

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=custom -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	fragmentPath := filepath.Join(ctx.Dir, "dev-dep", "Vagrantfile.fragment")
	return compile.App(ctx, &compile.AppOptions{
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context: map[string]interface{}{
				"fragment_path": fragmentPath,
			},
		},
		Customizations: []*compile.Customization{
			&compile.Customization{
				Type:     "dev-dep",
				Callback: processCustomDevDep,
				Schema: map[string]*schema.FieldSchema{
					"vagrantfile": &schema.FieldSchema{
						Type:        schema.TypeString,
						Description: "Path to Vagrantfile template",
					},
				},
			},
		},
	})
}

func (a *App) Build(ctx *app.Context) error {
	return nil
}

func (a *App) Deploy(ctx *app.Context) error {
	return nil
}

func (a *App) Dev(ctx *app.Context) error {
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
	}).Route(ctx)
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	return nil, nil
}

func processCustomDevDep(d *schema.FieldData) (*compile.CustomizationResult, error) {
	if _, ok := d.GetOk("vagrantfile"); !ok {
		return nil, fmt.Errorf(
			"Customization 'dev-dep': 'vagrantfile' must be specified")
	}

	return &compile.CustomizationResult{
		Callback: compileDev(d),
	}, nil
}

func compileDev(d *schema.FieldData) compile.CompileCallback {
	vf := d.Get("vagrantfile").(string)

	return func(ctx *app.Context, data *bindata.Data) error {
		fragment := data.Context["fragment_path"].(string)
		if err := data.RenderReal(fragment, vf); err != nil {
			return err
		}

		return data.RenderAsset(
			filepath.Join(ctx.Dir, "dev", "Vagrantfile"),
			"data/dev/Vagrantfile.tpl")
	}
}

const devInstructions = `
A development has been created.

Note that this development environment is just an example of what a
consumer of this application might see as a development dependency.
"Custom" types are not meant to be mutably developed like normal
applications.
`
