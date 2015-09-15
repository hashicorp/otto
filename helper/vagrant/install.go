package vagrant

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/hashicorp/go-checkpoint"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/otto/ui"
)

var (
	versionMin = version.Must(version.NewVersion("1.7.4"))
	versionRe  = regexp.MustCompile(`Vagrant (.+)`)
)

// InstallIfNeeded will install Vagrant if needed, but will prompt the user
// and output various status messages to the given UI.
func InstallIfNeeded(u ui.Ui) error {
	// First check if Vagrant is installed
	vsn, err := Version()
	if err != nil {
		return err
	}

	// Start building the input if we need to ask to install Vagrant,
	// but check if we perhaps have an acceptable version already.
	input := &ui.InputOpts{
		Id:      "vagrant_install",
		Query:   "Would you like Otto to install Vagrant?",
		Default: "",
	}
	if vsn != nil {
		log.Printf("[DEBUG] vagrant installation found: %s", vsn)
		if !vsn.LessThan(versionMin) {
			log.Printf(
				"[DEBUG] vagrant installation not necessary, min: %s",
				versionMin)
			return nil
		}

		log.Printf(
			"[INFO] vagrant installation required, min not reached: %s",
			versionMin)
		input.Description = fmt.Sprintf(
			vagrantInstallUpdateRequired, vsn, versionMin)
	} else {
		// We don't have a version, so we don't have Vagrant installed
		log.Printf("[INFO] vagrant installation required, vagrant not found")
		input.Description = vagrantInstallRequired
	}

	// Ask the user
	result, err := u.Input(input)
	if err != nil {
		return err
	}
	if strings.ToLower(result) != "yes" {
		return fmt.Errorf("Installation cancelled")
	}

	// We want to install
	return Install(u)
}

// Install will install the latest version of Vagrant.
func Install(u ui.Ui) error {
	// Determine the latest version using Checkpoint
	u.Header("Determining latest version of Vagrant to install...")
	check, err := checkpoint.Check(&checkpoint.CheckParams{
		Product: "vagrant",
		Force:   true,
	})
	if err != nil {
		return err
	}
	u.Message(fmt.Sprintf("Latest version: %s", check.CurrentVersion))

	// Determine the URL
	var url string
	switch runtime.GOOS {
	case "darwin":
		url = fmt.Sprintf(
			"https://dl.bintray.com/mitchellh/vagrant/vagrant_%s.dmg",
			check.CurrentVersion)
		return installDarwin(u, url)
	default:
		return fmt.Errorf(
			"Otto doesn't yet support installing Vagrant automatically\n" +
				"on your OS. Please install Vagrant manually.")
	}
}

// Version returns the installed version of Vagrant.
//
// If Vagrant is detected as not installed, an empty version string will
// be returned with no error. An error is only returned if there is an error
// detecting Vagrant or its version.
func Version() (*version.Version, error) {
	// Look for it on the path first
	_, err := exec.LookPath("vagrant")
	if err != nil {
		if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
			return nil, nil
		}

		return nil, err
	}

	// Grab the version
	var buf bytes.Buffer
	cmd := exec.Command("vagrant", "--version")
	cmd.Stdout = &buf
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Match the version out
	matches := versionRe.FindStringSubmatch(buf.String())
	if len(matches) == 0 {
		return nil, fmt.Errorf("unable to find Vagrant version: %s", buf.String())
	}

	return version.NewVersion(matches[1])
}

// installDarwin installs on OS X
func installDarwin(u ui.Ui, url string) error {
	u.Header("Installing Vagrant for OS X...")
	u.Message(fmt.Sprintf("Installer URL: %s", url))
	u.Message("")
	u.Message(
		"Otto will now execute a local script to download and install\n" +
			"Vagrant from the address above. At some point in this script\n" +
			"administrator privileges will be required and your password\n" +
			"will be asked for.\n\n" +
			"If you're uncomfortable with Otto doing this automatically, you\n" +
			"can interrupt the process with Ctrl-C at any point and install\n" +
			"Vagrant manually.\n\n")

	// Grab the script from our assets
	asset, err := Asset("data/darwin/install.sh")
	if err != nil {
		return err
	}

	// Create a temporary directory for our script
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		return err
	}
	defer os.RemoveAll(td)

	// Write the script
	path := filepath.Join(td, "install.sh")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(asset))
	f.Close()
	if err != nil {
		return err
	}

	// Execute the script
	cmd := exec.Command("bash", path, url)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Success!
	u.Header("[green]Vagrant installed!")
	return nil
}

const vagrantInstallUpdateRequired = `
An older version of Vagrant was found installed (%s). Otto requires
version %s or higher. Otto can install the latest version of Vagrant for
you. Would you like Otto to install Vagrant for you? Vagrant will be
installed system-wide using official Vagrant installers.

Please enter 'yes' to continue. Any other value will exit.
`

const vagrantInstallRequired = `
Otto uses software called Vagrant to manage development environments.
Vagrant couldn't be found on your system. Otto can install the latest
version of Vagrant for you automatically. Would you like Otto to install
Vagrant for you?

Please enter 'yes' to continue. Any other value will exit.
`
