package vagrant

import (
	"fmt"

	"github.com/hashicorp/otto/app"
)

type DevDepOptions struct {
	// Dir is the directory where Vagrant will be executed.
	Dir string

	// Script is the script to run to build the dev dependency.
	Script string

	// Files are the resulting files relative to the cache directory
	// that are part of the dep. If these don't exist, an error will be
	// generated.
	Files []string
}

// DevDep builds a dev dependency using Vagrant.
//
// This function uses Build to build the dev dependency. Please see
// the documentation of that function for more details on how that works.
//
// This function implements app.App.DevDep.
func DevDep(dst *app.Context, src *app.Context, opts *DevDepOptions) (*app.DevDep, error) {
	src.Ui.Header(fmt.Sprintf(
		"Building the dev dependency: '%s'", src.Appfile.Application.Name))
	src.Ui.Message(
		"To ensure cross-platform compatibility, we'll use Vagrant to\n" +
			"build this application. This is slow, and in a lot of cases we\n" +
			"can do something faster. Future versions of Otto will detect and\n" +
			"do this. As long as the application doesn't change, Otto will\n" +
			"cache the results of this build.\n\n")

	// Use the Build function to do so...
	err := Build(src, &BuildOptions{
		Dir:    opts.Dir,
		Script: opts.Script,
	})
	if err != nil {
		return nil, err
	}

	// Return the dep with the configured files. Eventually we'll verify
	// these files exist. For now, we don't.
	return &app.DevDep{Files: opts.Files}, nil
}
