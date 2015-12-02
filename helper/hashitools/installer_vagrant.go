package hashitools

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/otto/ui"
	"github.com/mitchellh/ioprogress"
	"gopkg.in/flosch/pongo2.v3"
)

// VagrantInstaller is an Installer that knows how to install Vagrant,
// which uses its own system installer.
type VagrantInstaller struct {
	// Ui is the Otto UI for asking the user for input and outputting
	// the status of installation.
	Ui ui.Ui
}

func (i *VagrantInstaller) InstallAsk(installed, required, latest *version.Version) (bool, error) {
	input := &ui.InputOpts{
		Id:      "vagrant_install",
		Query:   "Would you like Otto to install Vagrant?",
		Default: "",
	}

	// Figure out the description text to use for input
	var tplString string
	if installed == nil {
		tplString = vagrantInstallRequired
	} else {
		tplString = vagrantInstallRequiredUpdate
	}

	// Parse the template and render it
	tpl, err := pongo2.FromString(strings.TrimSpace(tplString))
	if err != nil {
		return false, err
	}
	input.Description, err = tpl.Execute(map[string]interface{}{
		"name":      "vagrant",
		"installed": installed,
		"latest":    latest,
		"required":  required,
	})
	if err != nil {
		return false, err
	}

	result, err := i.Ui.Input(input)
	if err != nil {
		return false, err
	}

	return strings.ToLower(result) == "yes", nil
}

func (i *VagrantInstaller) Install(vsn *version.Version) error {
	// Determine the URL
	var url string
	switch runtime.GOOS {
	case "darwin":
		url = fmt.Sprintf(
			"https://releases.hashicorp.com/vagrant/%s/vagrant_%s.dmg",
			vsn, vsn)
	default:
		return fmt.Errorf(
			"Otto doesn't yet support installing Vagrant automatically\n" +
				"on your OS. Please install Vagrant manually.")
	}

	// Create the temporary directory where we'll store the data
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		return err
	}
	defer os.RemoveAll(td)

	// Create the file path
	path := filepath.Join(td, "vagrant")
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	// Download the file to a temporary directory
	i.Ui.Header(fmt.Sprintf("Downloading Vagrant v%s...", vsn))
	i.Ui.Message("URL: " + url)
	i.Ui.Message("")
	resp, err := cleanhttp.DefaultClient().Get(url)
	if err != nil {
		f.Close()
		return err
	}
	if resp.StatusCode != 200 {
		resp.Body.Close()
		f.Close()
		return fmt.Errorf("Error downloading, status code %d", resp.StatusCode)
	}

	// Build the progress bar for our download
	progressR := &ioprogress.Reader{
		Reader:   resp.Body,
		Size:     resp.ContentLength,
		DrawFunc: ioprogress.DrawTerminalf(os.Stdout, i.progressFormat),
	}

	// Listen for interrupts so we can cancel the download
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	// Copy the zip data
	errCh := make(chan error, 1)
	go func() {
		_, err := io.Copy(f, progressR)
		errCh <- err
	}()

	// Wait for an interrupt or finish
	select {
	case err = <-errCh:
	case <-sigCh:
		err = fmt.Errorf("interrupted")
	}

	// Finish up
	resp.Body.Close()
	f.Close()
	if err != nil {
		return err
	}

	// Run the proper installer
	switch runtime.GOOS {
	case "darwin":
		return i.installDarwin(path)
	}

	panic("unknown OS")
}

func (i *VagrantInstaller) Path() string {
	// Always return "" since "vagrant" is just expected to be on the PATH
	return ""
}

func (i *VagrantInstaller) installDarwin(installerPath string) error {
	// Grab the script from our assets
	asset, err := Asset("data/vagrant-darwin/install.sh")
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
	cmd := exec.Command("bash", path, installerPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	i.Ui.Header("[green]Vagrant installed successfully!")
	return nil
}

func (i *VagrantInstaller) progressFormat(progress, total int64) string {
	return fmt.Sprintf("    %s", ioprogress.DrawTextFormatBytes(progress, total))
}

const vagrantInstallRequired = `
Otto requires {{name}} to be installed, but it couldn't be found on your
system. Otto can install the latest version of {{name}} for you. {{name}}
will be installed system-wide since it uses system-specific installers.
Would you like Otto to install {{name}} for you? Alternatively, you may install
this on your own.

If you answer yes, Otto will install {{name}} version {{latest}}.

Please enter 'yes' to continue. Any other value will exit.
`

const vagrantInstallRequiredUpdate = `
An older version of {{name}} was found installed ({{installed}}). Otto requires
version {{required}} or higher. Otto can install the latest version of {{name}}
for you ({{latest}}). Would you like Otto to update {{name}} for you?

Please enter 'yes' to continue. Any other value will exit.
`
