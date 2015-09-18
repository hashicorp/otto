package vagrant

import (
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/otto/ui"
)

// SSHCache is a helper to cache the SSH connection info from Vagrant
// and use that for executing to avoid the overhead of loading Vagrant.
type SSHCache struct {
	// Path is the path to where the SSH cache file should go
	Path string

	// Vagrant is the Vagrant instance we'll use to execute Vagrant commands.
	Vagrant *Vagrant
}

// Exec executes SSH and opens a console.
//
// This will use the cached SSH info if it exists, or will otherwise
// drop into `vagrant ssh`. If cacheOkay is false, then it'll always go
// straight to `vagrant ssh`.
func (c *SSHCache) Exec(cacheOkay bool) error {
	// If we have the cache file, use that
	if _, err := os.Stat(c.Path); err == nil {
		cmd := exec.Command("ssh", "-F", c.Path, "default")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil {
			return err
		}
		return nil
	}

	// Otherwise raw SSH
	return c.Vagrant.Execute("ssh")
}

// Cache will execute "ssh-config" and cache the SSH info.
func (c *SSHCache) Cache() error {
	// We just copy the Vagrant instance so we can modify it without
	// worrying about restoring stuff.
	var mockUi ui.Mock
	vagrant := *c.Vagrant
	vagrant.Ui = &mockUi
	if err := vagrant.Execute("ssh-config"); err != nil {
		return err
	}

	// Write the output to the cache
	f, err := os.Create(c.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, raw := range mockUi.RawBuf {
		if _, err := io.Copy(f, strings.NewReader(raw)); err != nil {
			return err
		}
	}

	return nil
}

// Delete clears the cache.
func (c *SSHCache) Delete() error {
	// We ignore the return value here because it'll happen if the
	// file doesn't exist and we just don't care.
	os.Remove(c.Path)
	return nil
}
