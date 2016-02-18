package otto

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
)

// infraCreds configures the credentials on the context. If we don't
// have credentials, then this will handle querying them from the
// user. If we do have the credentials, then this will decrypt and
// read them.
func (c *Core) infraCreds(infra infrastructure.Infrastructure, ctx *infrastructure.Context) error {
	// Output to the user some information about what is about to happen here...
	ctx.Ui.Header(fmt.Sprintf(
		"Detecting infrastructure credentials for: %s (%s)",
		ctx.Infra.Name, ctx.Infra.Type))

	// The path to where we put the encrypted creds. This is always local.
	// We don't do remote storage of credentials. That is up to the user
	// to do if wanted (via something like Vault).
	path := c.infraCredsPath(ctx)

	// Determine whether we believe the creds exist already or not
	var exists bool
	if _, err := os.Stat(path); err == nil {
		exists = true
	}

	// If the creds exist, then attempt to decrypt them, asking the
	// user for the password.
	var creds map[string]string
	if exists {
		ctx.Ui.Message(
			"Cached and encrypted infrastructure credentials found.\n" +
				"Otto will now ask you for the password to decrypt these\n" +
				"credentials.\n\n")

		// If they exist, ask for the password
		value, err := ctx.Ui.Input(&ui.InputOpts{
			Id:          "creds_password",
			Query:       "Encrypted Credentials Password",
			Description: strings.TrimSpace(credsQueryPassExists),
			Hide:        true,
			EnvVars:     []string{"OTTO_CREDS_PASSWORD"},
		})
		if err != nil {
			return err
		}

		// If the password is not blank, then just read the credentials
		if value != "" {
			plaintext, err := cryptRead(path, value)
			if err == nil {
				err = json.Unmarshal(plaintext, &creds)
			}
			if err != nil {
				return fmt.Errorf(
					"error reading encrypted credentials: %s\n\n"+
						"If this error persists, you can force Otto to ask for credentials\n"+
						"again by inputting the empty password as the password.",
					err)
			}
		}
	}

	// If we don't have creds, then we need to query the user via the
	// infrastructure implementation. This could be because we don't have
	// them saved or the user explicitly requested it.
	if creds == nil {
		ctx.Ui.Message(
			"Existing infrastructure credentials were not found! Otto will\n" +
				"now ask you for infrastructure credentials. These will be encrypted\n" +
				"and saved on disk so this doesn't need to be repeated.\n\n" +
				"IMPORTANT: If you're re-entering new credentials, make sure the\n" +
				"credentials are for the same account, otherwise you may lose\n" +
				"access to your existing infrastructure Otto set up.\n\n")

		// Ask the infra implementation for them
		var err error
		creds, err = infra.Creds(ctx)
		if err != nil {
			return err
		}

		// Now that we have the credentials, we need to ask for the
		// password to encrypt and store them.
		var password string
		for password == "" {
			password, err = ctx.Ui.Input(&ui.InputOpts{
				Id:          "creds_password",
				Query:       "Password for Encrypting Credentials",
				Description: strings.TrimSpace(credsQueryPassNew),
				Hide:        true,
				EnvVars:     []string{"OTTO_CREDS_PASSWORD"},
			})
			if err != nil {
				return err
			}
		}

		// With the password, encrypt and write the data
		plaintext, err := json.Marshal(creds)
		if err != nil {
			// creds is a map[string]string, so this shouldn't ever fail
			panic(err)
		}

		// Create the directory
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf(
				"error writing encrypted credentials: %s", err)
		}

		if err := cryptWrite(path, password, plaintext); err != nil {
			return fmt.Errorf(
				"error writing encrypted credentials: %s", err)
		}
	}

	// Set the credentials
	ctx.InfraCreds = creds

	// Let the infrastructure do whatever it likes to verify that the credentials
	// are good, so we can fail fast in case there's a problem.
	return infra.VerifyCreds(ctx)
}

// infraCredsPath returns the path to the encrypted infrastructure
// credentials if we have those.
func (c *Core) infraCredsPath(ctx *infrastructure.Context) string {
	return filepath.Join(c.dataDir, "cache", "creds", ctx.Infra.Name)
}

const credsQueryPassExists = `
Infrastructure credentials are required for this operation. Otto found
saved credentials that are password protected. Please enter the password
to decrypt these credentials. You may also just hit <enter> and leave
the password blank to force Otto to ask for the credentials again.
`

const credsQueryPassNew = `
This password will be used to encrypt and save the credentials so they
don't need to be repeated multiple times.
`
