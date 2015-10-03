package hashitools

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/go-checkpoint"
	"github.com/hashicorp/go-version"
)

//go:generate go-bindata -pkg=hashitools -nomemcopy -nometadata ./data/...

var (
	versionRe = regexp.MustCompile(`v?(\d+\.\d+\.[^\s]+)`)
)

// Project represents a HashiCorp Go project and provides various operations
// around that.
type Project struct {
	// Name is the name of the project, all lowercase
	Name string

	// Installer is the installer for this project
	Installer Installer

	// MinVersion is the minimum version of this project that Otto
	// can use to function. This will be used with `InstallIfNeeded`
	// to prompt the user to install.
	MinVersion *version.Version
}

// InstallIfNeeded will check if installation of this project is required
// and will invoke the installer if needed.
func (p *Project) InstallIfNeeded() error {
	log.Printf("[DEBUG] installIfNeeded: %s", p.Name)

	// Start grabbing the latest version as early as possible since
	// this requires a network call. We might as well do it while we're
	// doing a subprocess.
	latestCh := make(chan *version.Version, 1)
	errCh := make(chan error, 1)
	go func() {
		latest, err := p.LatestVersion()
		if err != nil {
			errCh <- err
		}
		latestCh <- latest
	}()

	// Grab the version we have installed
	installed, err := p.Version()
	if err != nil {
		return err
	}
	if installed == nil {
		log.Printf("[DEBUG] installIfNeeded: %s not installed", p.Name)
	} else {
		log.Printf("[DEBUG] installIfNeeded: %s installed: %s", p.Name, installed)
	}

	// Wait for the latest
	var latest *version.Version
	select {
	case latest = <-latestCh:
	case err := <-errCh:
		return err
	}
	log.Printf("[DEBUG] installIfNeeded: %s latest: %s", p.Name, latest)
	log.Printf("[DEBUG] installIfNeeded: %s min: %s", p.Name, p.MinVersion)

	// Determine if we require an install
	installRequired := installed == nil
	if installed != nil {
		if installed.LessThan(p.MinVersion) {
			installRequired = true
		}

		// TODO: updates
	}

	// No install required? Exit out.
	if !installRequired {
		log.Printf("[DEBUG] installIfNeeded: %s no installation needed", p.Name)
		return nil
	}

	// We need to install! Ask the installer to verify for us
	ok, err := p.Installer.InstallAsk(installed, p.MinVersion, latest)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("Installation cancelled")
	}

	// Install
	return p.Installer.Install(latest)
}

// Latest version returns the latest version of this project.
func (p *Project) LatestVersion() (*version.Version, error) {
	check, err := checkpoint.Check(&checkpoint.CheckParams{
		Product: p.Name,
		Force:   true,
	})
	if err != nil {
		return nil, err
	}

	return version.NewVersion(check.CurrentVersion)
}

// Version reads the version of this project.
func (p *Project) Version() (*version.Version, error) {
	path := p.Path()
	if !filepath.IsAbs(path) {
		// Look for it on the path first if we don't have a full path to it
		_, err := exec.LookPath(path)
		if err != nil {
			if execErr, ok := err.(*exec.Error); ok && execErr.Err == exec.ErrNotFound {
				return nil, nil
			}

			return nil, err
		}
	}

	// Grab the version
	var stdout, buf bytes.Buffer
	cmd := exec.Command(path, "--version")
	cmd.Stdout = io.MultiWriter(&stdout, &buf)
	cmd.Stderr = &buf
	runErr := cmd.Run()

	// Match the version out before we check for a run error, since some `project
	// --version` commands can return a non-zero exit code.
	matches := versionRe.FindStringSubmatch(stdout.String())
	if len(matches) == 0 {
		if runErr != nil {
			return nil, fmt.Errorf(
				"Error checking %s version: %s\n\n%s",
				p.Name, runErr, buf.String())
		}
		return nil, fmt.Errorf(
			"unable to find %s version in output: %q", p.Name, buf.String())
	}

	return version.NewVersion(matches[1])
}

// Path returns the path to this project. This will check if the project
// binary is pre-installed in our installation directory and use that path.
// Otherwise, it will return the raw project name.
func (p *Project) Path() string {
	if p := p.Installer.Path(); p != "" {
		return p
	}

	return p.Name
}
