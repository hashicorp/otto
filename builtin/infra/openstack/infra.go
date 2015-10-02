package openstack

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/sshagent"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
	"github.com/mitchellh/go-homedir"
)

//go:generate go-bindata -pkg=openstack -nomemcopy -nometadata ./data/...

// Infra returns the infrastructure.Infrastructure implementation.
// This function is a infrastructure.Factory.
func Infra() (infrastructure.Infrastructure, error) {
	return &terraform.Infrastructure{
		CredsFunc:       creds,
		VerifyCredsFunc: verifyCreds,
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
		},
		Variables: map[string]string{},
	}, nil
}

func creds(ctx *infrastructure.Context) (map[string]string, error) {
	fields := []*ui.InputOpts{
		&ui.InputOpts{
			Id:          "openstack_auth_url",
			Query:       "OpenStack authentification URL endpoint",
			Description: "OpenStack authentification URL endpoint (Keystone v2 only)",
			EnvVars:     []string{"OS_AUTH_URL"},
		},
		&ui.InputOpts{
			Id:          "openstack_tenant_name",
			Query:       "OpenStack tenant name",
			Description: "OpenStack tenant name",
			EnvVars:     []string{"OS_TENANT_NAME"},
		},
		&ui.InputOpts{
			Id:          "openstack_username",
			Query:       "OpenStack user name",
			Description: "OpenStack user name",
			EnvVars:     []string{"OS_USERNAME"},
		},
		&ui.InputOpts{
			Id:          "openstack_password",
			Query:       "OpenStack user password",
			Description: "OpenStack user password",
			EnvVars:     []string{"OS_PASSWORD"},
		},
		&ui.InputOpts{
			Id:          "openstack_region_name",
			Query:       "OpenStack region name",
			Description: "OpenStack region name",
			EnvVars:     []string{"OS_REGION_NAME"},
		},
		&ui.InputOpts{
			Id:          "openstack_image_id",
			Query:       "OpenStack image ID",
			Description: "Image ID for the OpenStack instances",
			EnvVars:     []string{"OS_IMAGE_ID"},
		},
		&ui.InputOpts{
			Id:          "openstack_flavor_id",
			Query:       "OpenStack flavor ID",
			Description: "Flavor ID for the OpenStack instances",
			EnvVars:     []string{"OS_FLAVOR_ID"},
		},
		&ui.InputOpts{
			Id:          "openstack_floating_ip_pool",
			Query:       "OpenStack floating IP pool",
			Description: "OpenStack floating IP pool used to get public IPs",
			EnvVars:     []string{"OS_FLOATING_IP_POOL"},
		},
		&ui.InputOpts{
			Id:          "openstack_ssh_username",
			Query:       "OpenStack SSH username",
			Description: "Username used to SSH into the OpenStack instances",
			EnvVars:     []string{"OS_SSH_USERNAME"},
		},
		&ui.InputOpts{
			Id:          "openstack_external_gateway",
			Query:       "OpenStack extrernal gateway",
			Description: "External gateway that will be used to connect routers",
			EnvVars:     []string{"OS_EXTERNAL_GATEWAY_ID"},
		},
		&ui.InputOpts{
			Id:          "ssh_public_key_path",
			Query:       "SSH Public Key Path",
			Description: "Path to an SSH public key that will be granted access to OpenStack instances",
			Default:     "~/.ssh/id_rsa.pub",
			EnvVars:     []string{"OS_SSH_PUBLIC_KEY_PATH"},
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

func verifyCreds(ctx *infrastructure.Context) error {
	found, err := sshagent.HasKey(ctx.InfraCreds["ssh_public_key"])
	if err != nil {
		return sshAgentError(err)
	}
	if !found {
		ok, _ := guessAndLoadPrivateKey(
			ctx.Ui, ctx.InfraCreds["ssh_public_key_path"])
		if ok {
			ctx.Ui.Message(
				"A private key was found and loaded. Otto will now check\n" +
					"the SSH Agent again and continue if the correct key is loaded")

			found, err = sshagent.HasKey(ctx.InfraCreds["ssh_public_key"])
			if err != nil {
				return sshAgentError(err)
			}
		}
	}

	if !found {
		return sshAgentError(fmt.Errorf(
			"You specified an SSH public key of: %q, but the private key from this\n"+
				"keypair is not loaded the SSH Agent. To load it, run:\n\n"+
				"  ssh-add [PATH_TO_PRIVATE_KEY]",
			ctx.InfraCreds["ssh_public_key_path"]))
	}
	return nil
}

func sshAgentError(err error) error {
	return fmt.Errorf(
		"Otto uses your SSH Agent to authenticate with instances created in\n"+
			"AWS, but it could not verify that your SSH key is loaded into the agent.\n"+
			"The error message follows:\n\n%s", err)
}

// guessAndLoadPrivateKey takes a path to a public key and determines if a
// private key exists by just stripping ".pub" from the end of it. if so,
// it attempts to load that key into the agent.
func guessAndLoadPrivateKey(ui ui.Ui, pubKeyPath string) (bool, error) {
	fullPath, err := homedir.Expand(pubKeyPath)
	if err != nil {
		return false, err
	}
	if !strings.HasSuffix(fullPath, ".pub") {
		return false, fmt.Errorf("No .pub suffix, cannot guess path.")
	}
	privKeyGuess := strings.TrimSuffix(fullPath, ".pub")
	if _, err := os.Stat(privKeyGuess); os.IsNotExist(err) {
		return false, fmt.Errorf("No file at guessed path.")
	}

	ui.Header("Loading key into SSH Agent")
	ui.Message(fmt.Sprintf(
		"The key you provided (%s) was not found in your SSH Agent.", pubKeyPath))
	ui.Message(fmt.Sprintf(
		"However, Otto found a private key here: %s", privKeyGuess))
	ui.Message(fmt.Sprintf(
		"Automatically running 'ssh-add %s'.", privKeyGuess))
	ui.Message("If your SSH key has a passphrase, you will be prompted for it.")
	ui.Message("")

	if err := sshagent.Add(ui, privKeyGuess); err != nil {
		return false, err
	}

	return true, nil
}
