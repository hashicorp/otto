package vagrant

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/hashicorp/go-version"
)

var (
	versionRe = regexp.MustCompile(`^Vagrant (.+?)$`)
)

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
