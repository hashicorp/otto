package hashigo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/go-checkpoint"
	"github.com/hashicorp/go-version"
)

var (
	versionRe = regexp.MustCompile(`v?(\d+\.\d+\.[^\s]+)`)
)

// Project represents a HashiCorp Go project and provides various operations
// around that.
type Project struct {
	// Name is the name of the project, all lowercase
	Name string

	// InstallDir is the directory where we will install the project.
	// This will only work for Go projects distributed as zips.
	InstallDir string

	// MinVersion is the minimum version of this project that Otto
	// can use to function. This will be used with `InstallIfNeeded`
	// to prompt the user to install.
	MinVersion *version.Version
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
	var buf bytes.Buffer
	cmd := exec.Command(path, "--version")
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

// Path returns the path to this project. This will check if the project
// binary is pre-installed in our installation directory and use that path.
// Otherwise, it will return the raw project name.
func (p *Project) Path() string {
	full := filepath.Join(p.InstallDir, p.Name, p.Name)
	if _, err := os.Stat(full); err == nil {
		return full
	}

	return p.Name
}
