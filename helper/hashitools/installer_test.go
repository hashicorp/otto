package hashitools

import (
	"testing"
)

func TestGoInstaller_impl(t *testing.T) {
	var _ Installer = new(GoInstaller)
}
