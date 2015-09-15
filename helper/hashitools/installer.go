package hashitools

import (
	"github.com/hashicorp/go-version"
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
