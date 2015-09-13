package aws

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
	"github.com/mitchellh/go-homedir"
)

//go:generate go-bindata -pkg=aws -nomemcopy -nometadata ./data/...

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
