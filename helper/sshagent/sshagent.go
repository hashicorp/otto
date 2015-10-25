// Helpers for interacting with the local SSH Agent
package sshagent

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"reflect"

	execHelper "github.com/hashicorp/otto/helper/exec"
	"github.com/hashicorp/otto/ui"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// HasKey determines if a given public key (provided as a string with the
// contents of a public key file), is loaded into the local SSH Agent.
func HasKey(publicKey string) (bool, error) {
	pk, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return false, fmt.Errorf("Error parsing provided public key: %s", err)
	}

	agentKeys, err := ListKeys()
	if err != nil {
		return false, err
	}

	for _, agentKey := range agentKeys {
		if reflect.DeepEqual(agentKey.Marshal(), pk.Marshal()) {
			return true, nil
		}
	}
	return false, nil
}

// ListKeys connects to the local SSH Agent and lists all the public keys
// loaded into it. It returns user friendly error message when it has trouble.
func ListKeys() ([]*agent.Key, error) {
	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		return nil, fmt.Errorf(
			"The SSH_AUTH_SOCK environment variable is not set, which normally\n" +
				"means that no SSH Agent is running.")
	}

	conn, err := net.Dial("unix", sshAuthSock)
	if err != nil {
		return nil, fmt.Errorf(
			"Error connecting to agent: %s\n\n"+
				"The agent address is detected using the SSH_AUTH_SOCK environment\n"+
				"variable. Please verify this variable is correct and the SSH agent\n"+
				"is properly set up.",
			err)
	}
	defer conn.Close()

	agent := agent.NewClient(conn)
	loadedKeys, err := agent.List()
	if err != nil {
		return nil, fmt.Errorf("Error listing keys: %s", err)
	}
	return loadedKeys, err
}

// Add takes the path of a private key and runs ssh-add locally to add it to
// the agent. It needs a Ui to be able to interact with the user for the
// password prompt.
func Add(ui ui.Ui, privateKeyPath string) error {
	cmd := exec.Command("ssh-add", privateKeyPath)
	return execHelper.Run(ui, cmd)
}
