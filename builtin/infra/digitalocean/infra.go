package digitalocean

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

//go:generate go-bindata -pkg=digitalocean -nomemcopy -nometadata ./data/...

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
		Variables: map[string]string{
			"do_region": "sfo1",
		},
	}, nil
}

func creds(ctx *infrastructure.Context) (map[string]string, error) {
	fields := []*ui.InputOpts{
		&ui.InputOpts{
			Id:          "do_token",
			Query:       "DigitalOcean API Token",
			Description: "DigitalOcean API Token",
			EnvVars:     []string{"DIGITALOCEAN_TOKEN"},
		},
		&ui.InputOpts{
			Id:          "do_ssh_public_key_path",
			Query:       "SSH Public Key Path",
			Description: "Path to an SSH public key that will be granted access to droplets",
			Default:     "~/.ssh/id_rsa.pub",
			EnvVars:     []string{"TF_DO_SSH_PUBLIC_KEY_PATH"},
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
	sshPath, err := homedir.Expand(result["do_ssh_public_key_path"])
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
			"DigitalOcean, but it could not verify that your SSH key is loaded into the agent.\n"+
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
