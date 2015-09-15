package hashitools

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/otto/ui"
	"github.com/mitchellh/ioprogress"
	"gopkg.in/flosch/pongo2.v3"
)

// Installer is the interface that knows how to install things.
//
// This is an interface to support different installation methods between
// our different projects.
type Installer interface {
	// InstallAsk should ask the user if they'd like to install the
	// project. This is only called if installation is actually required.
	InstallAsk(installed, required, latest *version.Version) (bool, error)

	// Install should install the specified version.
	Install(*version.Version) error

	// Path is the path to the installed main binary of this project,
	// or "" if it doesn't seem installed.
	Path() string
}

// GoInstaller is an Installer that knows how to install Go projects.
type GoInstaller struct {
	// Name is the name of the project to install
	Name string

	// Dir is the directory where projects will be installed. They will
	// be installed to a sub-directory of the project name. Example:
	// if Dir is "/foo", then the Packer binary would be installed to
	// "/foo/packer/packer"
	Dir string

	// Ui is the Otto UI for asking the user for input and outputting
	// the status of installation.
	Ui ui.Ui
}

func (i *GoInstaller) InstallAsk(installed, required, latest *version.Version) (bool, error) {
	input := &ui.InputOpts{
		Id:      fmt.Sprintf("%s_install", i.Name),
		Query:   fmt.Sprintf("Would you like Otto to install %s?", strings.Title(i.Name)),
		Default: "",
	}

	// Figure out the description text to use for input
	var tplString string
	if installed == nil {
		tplString = installRequired
	} else {
		tplString = installRequiredUpdate
	}

	// Parse the template and render it
	tpl, err := pongo2.FromString(strings.TrimSpace(tplString))
	if err != nil {
		return false, err
	}
	input.Description, err = tpl.Execute(map[string]interface{}{
		"name":      i.Name,
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

func (i *GoInstaller) Install(vsn *version.Version) error {
	// All Go projects use a standard URL format
	url := fmt.Sprintf(
		"https://dl.bintray.com/mitchellh/%s/%s_%s_%s_%s.zip",
		i.Name, i.Name, vsn, runtime.GOOS, runtime.GOARCH)

	// Create the temporary directory where we'll store the data
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		return err
	}
	defer os.RemoveAll(td)

	// Create the ZIP file
	zipPath := filepath.Join(td, "project.zip")
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}

	// Download the ZIP
	i.Ui.Header(fmt.Sprintf("Downloading %s v%s...", i.Name, vsn))
	i.Ui.Message("URL: " + url)
	i.Ui.Message("")
	resp, err := http.Get(url)
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

	// Open the zip file
	i.Ui.Header("Unzipping downloaded package...")
	zipR, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer zipR.Close()

	// Clear our install directory
	installDir := filepath.Join(i.Dir, i.Name)
	if err := os.RemoveAll(installDir); err != nil {
		return err
	}
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return err
	}

	// Copy all the files
	for _, f := range zipR.File {
		dst, err := os.OpenFile(
			filepath.Join(installDir, f.Name),
			os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
			f.Mode())
		if err != nil {
			return err
		}

		fr, err := f.Open()
		if err != nil {
			dst.Close()
			return err
		}

		_, err = io.Copy(dst, fr)
		fr.Close()
		dst.Close()
		if err != nil {
			return err
		}
	}

	i.Ui.Header(fmt.Sprintf("[green]%s installed successfully!", i.Name))
	return nil
}

func (i *GoInstaller) Path() string {
	path := filepath.Join(i.Dir, i.Name, i.Name)
	if _, err := os.Stat(path); err == nil {
		return path
	}

	return ""
}

func (i *GoInstaller) progressFormat(progress, total int64) string {
	return fmt.Sprintf("    %s", ioprogress.DrawTextFormatBytes(progress, total))
}

const installRequired = `
Otto requires {{name}} to be installed, but it couldn't be found on your
system. Otto can install the latest version of {{name}} for you. Otto will
install this into its own private data directory so it doesn't conflict
with anything else on your system. Would you like Otto to install {{name}}
for you? Alternatively, you may install this on your own.

If you answer yes, Otto will install {{name}} version {{latest}}.

Please enter 'yes' to continue. Any other value will exit.
`

const installRequiredUpdate = `
An older version of {{name}} was found installed ({{installed}}). Otto requires
version {{required}} or higher. Otto can install the latest version of {{name}}
for you ({{latest}}). Would you like Otto to update {{name}} for you?

Please enter 'yes' to continue. Any other value will exit.
`
