package aws

import (
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
)

//go:generate go-bindata -pkg=aws -nomemcopy ./data/...

// Infra returns the infrastructure.Infrastructure implementation.
// This function is a infrastructure.Factory.
func Infra() (infrastructure.Infrastructure, error) {
	return &terraform.Infrastructure{
		CredsFunc: creds,
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
		},
		Variables: map[string]string{
			"aws_region": "us-east-1",
		},
	}, nil
}

func creds(ctx *infrastructure.Context) (map[string]string, error) {
	fields := []*ui.InputOpts{
		&ui.InputOpts{
			Id:          "aws_access_key",
			Query:       "AWS Access Key",
			Description: "AWS access key used for API calls.",
		},
		&ui.InputOpts{
			Id:          "aws_secret_key",
			Query:       "AWS Secret Key",
			Description: "AWS secret key used for API calls.",
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

	return result, nil
}
