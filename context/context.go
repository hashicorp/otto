package context

import (
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/directory"
	"github.com/hashicorp/otto/ui"
)

// Shared is the shared contexts for app/infra.
type Shared struct {
	// InfraCreds are the credentials for working with the infrastructure.
	// These are guaranteed to be populated for the following function
	// calls:
	//
	//   App.Build
	//   TODO
	//
	InfraCreds map[string]string

	// Ui is the Ui object that can be used to communicate with the user.
	Ui ui.Ui

	// Directory is the directory service. This is available during
	// both execution and compilation and can be used to view the
	// global data prior to doing anything.
	Directory directory.Backend

	// InstallDir is the directory where binaries can be installed. Consider
	// this directory a cache: Otto may remove it at any point, although
	// unlikely. And you're responsible to clean up anything in here.
	InstallDir string

	// Appfile is the full appfile
	Appfile *appfile.File
}
