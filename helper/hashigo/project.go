package hashigo

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/go-version"
)

var (
	versionRe = regexp.MustCompile(`v?(\d+\.\d+\.[^\s]+)`)
)

// Project represents a HashiCorp Go project and provides various operations
// around that.
type Project struct {
	Name       string
	InstallDir string
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
