package dockerext

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
	"github.com/hashicorp/otto/helper/schema"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/helper/vagrant"
)

//go:generate go-bindata -pkg=dockerext -nomemcopy -nometadata ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Meta() (*app.Meta, error) {
	return Meta, nil
}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	fragmentPath := filepath.Join(ctx.Dir, "dev-dep", "Vagrantfile.fragment")

	var opts compile.AppOptions
	custom := &customizations{Opts: &opts}
	opts = compile.AppOptions{
		Ctx: ctx,
		FoundationConfig: foundation.Config{
			ServiceName: ctx.Application.Name,
		},
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
			Context: map[string]interface{}{
				"fragment_path": fragmentPath,
				"path": map[string]string{
					"guest_working": fmt.Sprintf(
						"/otto-deps/%s-%s",
						ctx.Application.Name,
						ctx.Appfile.ID),
				},
			},
		},
		Customizations: []*compile.Customization{
			&compile.Customization{
				Type:     "docker",
				Callback: custom.processDocker,
				Schema: map[string]*schema.FieldSchema{
					"image": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "",
						Description: "Image name to run",
					},

					"run_args": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "",
						Description: "Args to pass to `docker run`",
					},

					"command": &schema.FieldSchema{
						Type:        schema.TypeString,
						Default:     "",
						Description: "Command to pass to `docker run`",
					},
				},
			},
		},
	}

	return compile.App(&opts)
}

func (a *App) Build(ctx *app.Context) error {
	return nil
}

func (a *App) Deploy(ctx *app.Context) error {
	// Check if we have a deployment script
	path := filepath.Join(ctx.Dir, "deploy", "main.tf")
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf(deployError)
	}

	// We do! Run it
	return terraform.Deploy(&terraform.DeployOptions{
		DisableBuild: true,
		InfraOutputMap: map[string]string{
			"region":         "aws_region",
			"subnet-private": "private_subnet_id",
			"subnet-public":  "public_subnet_id",
		},
	}).Route(ctx)
}

func (a *App) Dev(ctx *app.Context) error {
	return vagrant.Dev(&vagrant.DevOptions{
		Instructions: strings.TrimSpace(devInstructions),
	}).Route(ctx)
}

func (a *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	// Nothing needs to be done for this
	return nil, nil
}

const devInstructions = `
A development environment has been created.

Note that this development environment is just an example of what a
consumer of this application might see as a development dependency.
"docker-external" application types represent Docker images that are
developed and built external to Otto. A future version of Otto will include
a native "docker" type for a Docker-based development workflow. For
the "docker-external" type, the specified docker image is started.
`

const deployError = `
Deployment isn't supported for "docker-external" yet.

This will be supported very soon. Otto plans to integrate with
Nomad (nomadproject.io) and once we do that, Otto will schedule the
container to run. Until then, deployment isn't supported. Sorry!
`
