package goapp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
)

// DetectImportPath will try to automatically determine the import path
// for the Go application under development.
//
// This is necessary to setup proper GOPATH directories for development
// and builds.
func DetectImportPath(ctx *app.Context) (string, error) {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		ctx.Ui.Message(
			"Warning! GOPATH not set. Otto will be unable to automatically\n" +
				"setup your application GOPATH for development and builds. While Otto\n" +
				"sets up a development for you, your folder structure outside of Otto\n" +
				"should still represent a proper Go environment. If you do this, then\n" +
				"the development and build process function a lot smoother.\n\n" +
				"For simple Go applications, this may not be necessary.\n\n" +
				"This is just an informational message. This is not a bug.")
		return "", nil
	}

	// Gopath should be absolute
	gopath, err := filepath.Abs(gopath)
	if err != nil {
		return "", fmt.Errorf(
			"Error expanding GOPATH to an absolute path: %s", err)
	}

	dir := filepath.Dir(ctx.Appfile.Path)
	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf(
			"Error expanding Appfile path to an absolute path: %s", err)
	}

	// If the directory to our Appfile is a symlink, resolve that symlink
	// through. This makes this heuristic work for local dependencies.
	if fi, err := os.Lstat(dir); err == nil {
		if fi.Mode()&os.ModeSymlink != 0 {
			newDir, err := os.Readlink(dir)
			if err != nil {
				return "", fmt.Errorf(
					"Error reading symlink %s: %s", dir, err)
			}

			dir = newDir
		}
	}

	// The directory has to be prefixed with the gopath
	gopath = filepath.Join(gopath, "src")
	if !strings.HasPrefix(dir, gopath) {
		ctx.Ui.Message(
			"Warning! It looks like your application is not within your set\n" +
				"GOPATH. Otto will be unable to automatically setup the proper\n" +
				"GOPATH structure within your development and build environments.\n\n" +
				"To fix this, please put your application into the proper GOPATH\n" +
				"location as according to standard Go development practices.")
		return "", nil
	}

	detected := dir[len(gopath)+1:]
	ctx.Ui.Message(fmt.Sprintf(
		"Detected import path: %s\n\n"+
			"Otto will use this import path to automatically setup your dev\n"+
			"and build environments in the proper directories.",
		detected))
	return detected, nil
}
