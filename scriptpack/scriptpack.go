// Package scriptpack is used to work with "ScriptPacks" which are
// packages of shell scripts that app types can define and use to get things
// done on the remote machine, whether its for development or deployment.
//
// ScriptPacks are 100% pure shell scripting. Any inputs must be received from
// environment variables. They aren't allowed to template at all. This is
// all done to ensure testability of the ScriptPacks.
//
// These are treated as first class elements within Otto to assist with
// testing.
//
// To create your own scriptpack, see the "template" folder within this
// directory. The folder structure and contents are important for scriptpacks
// to function correctly.
package scriptpack

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/hashicorp/atlas-go/archive"
	"github.com/hashicorp/otto/helper/bindata"
)

// ScriptPack is a struct representing a single ScriptPack. This is exported
// from the various scriptpacks.
type ScriptPack struct {
	// Name is an identifying name used to name the environment variables.
	// For example the root of where the files are will always be
	// SCRIPTPACK_<NAME>_ROOT.
	Name string

	// Data is the compiled bindata for this ScriptPack. The entire
	// AssetDirectory will be copied starting from "/data"
	Data bindata.Data

	// Dependencies are a list of other ScriptPacks that will always be
	// unpacked alongside this ScriptPack. By referencing actual ScriptPack
	// pointers, the dependencies will also be statically compiled into
	// the Go binaries that contain them.
	//
	// Dependencies can be accessed at the path specified by
	// SCRIPTPACK_<DEP>_ROOT. Note that if you depend on a ScriptPack
	// which itself has a conflicting named dependency, then the first
	// one loaded will win. Be careful about this.
	Dependencies []*ScriptPack
}

// Env returns the environment variables that should be set for this
// ScriptPack when it is executed.
//
// path is the path to the root of the directory where Write was called
// to write the ScriptPack output.
func (s *ScriptPack) Env(path string) map[string]string {
	result := make(map[string]string)
	result[fmt.Sprintf("SCRIPTPACK_%s_ROOT", s.Name)] = filepath.Join(path, s.Name)
	for _, dep := range s.Dependencies {
		for k, v := range dep.Env(path) {
			result[k] = v
		}
	}

	return result
}

// Write writes the contents of the ScriptPack and any dependencies into
// the given directory.
func (s *ScriptPack) Write(dst string) error {
	// Build the names of all scriptpacks
	spNames := make([]string, 1, len(s.Dependencies)+1)
	spNames[0] = s.Name
	for _, dep := range s.Dependencies {
		spNames = append(spNames, dep.Name)
	}

	// Deps
	for _, dep := range s.Dependencies {
		if err := dep.Write(dst); err != nil {
			return err
		}
	}

	// Our own
	if err := s.Data.CopyDir(filepath.Join(dst, s.Name), "data"); err != nil {
		return err
	}

	// Write the main file which has the env vars in it
	f, err := os.Create(filepath.Join(dst, "main.sh"))
	if err != nil {
		return err
	}
	defer f.Close()
	tpl := template.Must(template.New("root").Parse(mainShTpl))
	return tpl.Execute(f, map[string]interface{}{
		"scriptpacks": spNames,
	})
}

// WriteArchive writes the contents of the ScriptPack as a tar gzip to the
// given path.
func (s *ScriptPack) WriteArchive(dst string) error {
	// Let's just open the file we're going to write to first to verify
	// we can write there since everything else is pointless if we can't.
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a temporary directory to store the raw ScriptPack data
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		return err
	}
	defer os.RemoveAll(td)

	// Write the ScriptPack
	if err := s.Write(td); err != nil {
		return err
	}

	// Archive this ScriptPack
	a, err := archive.CreateArchive(td, &archive.ArchiveOpts{
		VCS: false,
	})
	if err != nil {
		return err
	}
	defer a.Close()

	// Write the archive to final path
	_, err = io.Copy(f, a)
	return err
}

const mainShTpl = `
#!/bin/bash

# Determine the directory of this script
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ]; do
  DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"
  SOURCE="$(readlink "$SOURCE")"
  [[ $SOURCE != /* ]] && SOURCE="$DIR/$SOURCE"
done
DIR="$( cd -P "$( dirname "$SOURCE" )" && pwd )"

# Set the env vars
{{range .scriptpacks}}
export SCRIPTPACK_{{ . }}_ROOT="${DIR}/{{ . }}"
{{end}}
`
