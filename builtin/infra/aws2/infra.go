package aws

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/plan"
	"github.com/hashicorp/otto/ui"
	"github.com/mitchellh/go-homedir"
)

//go:generate go-bindata -pkg=aws -nomemcopy -nometadata ./data/...

func Factory() (infrastructure.Infrastructure, error) {
	return &Infra{}, nil
}

// Infra implements the infrastructure.Infrastructure interface
type Infra struct{}

func (i *Infra) Creds(ctx *infrastructure.Context) (map[string]string, error) {
	fields := []*ui.InputOpts{
		&ui.InputOpts{
			Id:          "aws_access_key",
			Query:       "AWS Access Key",
			Description: "AWS access key used for API calls.",
			EnvVars:     []string{"AWS_ACCESS_KEY_ID"},
		},
		&ui.InputOpts{
			Id:          "aws_secret_key",
			Query:       "AWS Secret Key",
			Description: "AWS secret key used for API calls.",
			EnvVars:     []string{"AWS_SECRET_ACCESS_KEY"},
		},
		&ui.InputOpts{
			Id:          "ssh_public_key_path",
			Query:       "SSH Public Key Path",
			Description: "Path to an SSH public key that will be granted access to EC2 instances",
			Default:     "~/.ssh/id_rsa.pub",
			EnvVars:     []string{"AWS_SSH_PUBLIC_KEY_PATH"},
		},
	}

	result := make(map[string]string, len(fields))
	for _, f := range fields {
		value, err := ctx.Ui.Input(f)
		if err != nil {
			return nil, err
		}

		result[f.Id] = value
	}

	// Load SSH public key contents
	sshPath, err := homedir.Expand(result["ssh_public_key_path"])
	if err != nil {
		return nil, fmt.Errorf("Error expanding homedir for SSH key: %s", err)
	}

	sshKey, err := ioutil.ReadFile(sshPath)
	if err != nil {
		return nil, fmt.Errorf("Error reading SSH key: %s", err)
	}
	result["ssh_public_key"] = string(sshKey)

	return result, nil
}

func (i *Infra) VerifyCreds(*infrastructure.Context) error {
	return nil
}

func (i *Infra) Compile(ctx *infrastructure.Context) (*infrastructure.CompileResult, error) {
	data := &bindata.Data{
		Asset:    Asset,
		AssetDir: AssetDir,
	}

	// Just copy the data directory for our flavor as-is
	if err := data.CopyDir(ctx.Dir, "data/"+ctx.Infra.Flavor); err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *Infra) Plan(ctx *infrastructure.Context) ([]*plan.Plan, error) {
	// TODO: Eventually we'll inspect the version and decide based on that.
	// For now we just init at v0.
	lookup := &directory.InfraLookup{Name: ctx.Infra.Name}
	infra, err := ctx.Directory.GetInfra(lookup)
	if err != nil {
		return nil, err
	}
	if infra == nil {
		// Initialize the infrastructure
		infra = directory.NewInfra(ctx.Infra)
		if err := ctx.Directory.PutInfra(lookup, infra); err != nil {
			return nil, err
		}
	}

	// Parse the plans
	plans, err := plan.ParseFile(filepath.Join(ctx.Dir, "plans", "v0.hcl"))
	if err != nil {
		return nil, err
	}

	// Set common variables on all the plans
	inputs := map[string]interface{}{
		"context.compile_dir": ctx.Dir,
		"context.infra.name":  ctx.Infra.Name,
	}
	for _, p := range plans {
		if p.Inputs == nil {
			p.Inputs = make(map[string]interface{})
		}
		for k, v := range inputs {
			p.Inputs[k] = v
		}
	}

	return plans, nil

	// If no state, then we've never made the infra and we have to
	// If state, check for drift
	// Future: If deploy version changes, then move to that version

	// Or:
	// If no version, never made
	// If less than version, check for drift
	// If less than version, move to it
}
