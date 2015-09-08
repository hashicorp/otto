package goapp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/otto/app"
)

// detectImportPath will try to automatically determine the import path
// for the Go application under development.
//
// This is necessary to setup proper GOPATH directories for development
// and builds.
func detectImportPath(ctx *app.Context) (string, error) {
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

	gopath = filepath.Join(gopath, "src")
	dir := filepath.Dir(ctx.Appfile.Path)
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
