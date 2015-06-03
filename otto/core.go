package otto

import (
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/infrastructure"
)

// Core is the main struct to use to interact with Otto as a library.
type Core struct {
	appfile *appfile.File
	infras  map[string]infrastructure.Factory
}

// CoreConfig is configuration for creating a new core with NewCore.
type CoreConfig struct {
	// Appfile is the appfile that this core will be using for configuration.
	Appfile *appfile.File

	// Infrastructures is the map of available infrastructures. The
	// value is a factory that can create the infrastructure impl.
	Infrastructures map[string]infrastructure.Factory
}

// NewCore creates a new core.
//
// Once this function is called, this CoreConfig should not be used again
// or modified, since the Core may use parts of it without deep copying.
func NewCore(c *CoreConfig) (*Core, error) {
	return &Core{
		appfile: c.Appfile,
		infras:  c.Infrastructures,
	}, nil
}
