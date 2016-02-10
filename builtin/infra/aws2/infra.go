package aws

import (
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/plan"
)

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

func (i *Infra) Compile(*infrastructure.Context) (*infrastructure.CompileResult, error) {
	// If no state, then we've never made the infra and we have to
	// If state, check for drift
	// Future: If deploy version changes, then move to that version

	// Or:
	// If no version, never made
	// If less than version, check for drift
	// If less than version, move to it

	return nil, nil
}

func (i *Infra) Plan(*infrastructure.Context) ([]*plan.Plan, error) {
	return []*plan.Plan{
		&plan.Plan{
			Description: "Creating base infrastructure",
			Tasks: []*plan.Task{
				&plan.Task{
					Type:        "Terraform.Apply",
					Description: "Run Terraform to create the base infrastructure.",
					Args: map[string]*plan.TaskArg{
						"pwd": &plan.TaskArg{
							Value: "/foo",
						},
					},
				},

				&plan.Task{
					Type:        "Otto.PutInfra",
					Description: "Update the deployment version",
					Args: map[string]*plan.TaskArg{
						"deploy_version": &plan.TaskArg{
							Value: "1",
						},
					},
				},
			},
		},
	}, nil
}
