package aws

import (
	"path/filepath"

	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/plan"
)

//go:generate go-bindata -pkg=aws -nomemcopy -nometadata ./data/...

func Factory() (infrastructure.Infrastructure, error) {
	return &Infra{}, nil
}

// Infra implements the infrastructure.Infrastructure interface
type Infra struct{}

func (i *Infra) Creds(*infrastructure.Context) (map[string]string, error) {
	return nil, nil
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
	// Parse the plans
	plans, err := plan.ParseFile(filepath.Join(ctx.Dir, "plans", "v0.hcl"))
	if err != nil {
		return nil, err
	}

	// Set common variables on all the plans
	inputs := map[string]interface{}{
		"context.compile_dir": ctx.Dir,
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
