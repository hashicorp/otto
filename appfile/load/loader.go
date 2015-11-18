package load

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/appfile/detect"
	"github.com/hashicorp/otto/otto"
)

// Loader is used to load an Appfile.
//
// This logic used to live in the "compile" command directly but was
// extracted so it could be tested in isolation better and in case there
// is any interest in reusability.
//
// The purpose of the loader is to do the multi-step Appfile load outlined
// below. As input, we have the "real" Appfile, which is the real physical
// Appfile if it exists (or nil if it doesn't).
//
//   1.) Detect the type, known as the "detected" Appfile
//
//   2.) Merge with the detected Appfile with the real Appfile. If
//       the "detect" setting is set to "false", then we're done. Otherwise,
//       continue.
//
//   3.) Instantiate the proper plugin for the app type and call the Implicit
//       API to get an Appfile known as the "implicit" Appfile.
//
//   4.) Merge in the order: detected, implicit, real. Return the real
//       Appfile.
//
type Loader struct {
	// Detector is the detector configuration. If this is nil then
	// no type detection will be done.
	Detector *detect.Config

	// Compiler is the appfile compiler that we're using. This is used
	// to do a minimal compile (MinCompile) to realize imports of
	// Appfiles prior to implicit loading.
	Compiler *appfile.Compiler

	// Apps will be used to load the proper app implementation for
	// implicit loading.
	Apps map[app.Tuple]app.Factory
}

func (l *Loader) Load(f *appfile.File, dir string) (*appfile.File, error) {
	realFile := f

	// Load the "detected" Appfile
	appDef, err := appfile.Default(dir, l.Detector)
	if err != nil {
		return nil, fmt.Errorf("Error detecting app type: %s", err)
	}

	// Merge the detected Appfile with the real Appfile
	var merged appfile.File
	if err := merged.Merge(appDef); err != nil {
		return nil, fmt.Errorf("Error loading Appfile: %s", err)
	}
	if realFile != nil {
		if err := merged.Merge(realFile); err != nil {
			return nil, fmt.Errorf("Error loading Appfile: %s", err)
		}
	}
	realFile = &merged

	// If we have no application type, there is nothing more to do
	if realFile == nil || realFile.Application.Type == "" {
		return realFile, nil
	}

	// If we're not configured to do any further detection, return it
	if !realFile.Application.Detect {
		return realFile, nil
	}

	// Minimally compile the file that we can use to create a core
	compiled, err := l.Compiler.MinCompile(realFile)
	if err != nil {
		return nil, err
	}

	// Create a temporary directory we use for the core
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(td)

	// Create a core
	core, err := otto.NewCore(&otto.CoreConfig{
		DataDir:    filepath.Join(td, "data"),
		LocalDir:   filepath.Join(td, "local"),
		CompileDir: filepath.Join(td, "compile"),
		Appfile:    compiled,
		Apps:       l.Apps,
	})
	if err != nil {
		return nil, err
	}

	// Get the app implementation
	appImpl, appCtx, err := core.App()
	if err != nil {
		return nil, err
	}
	defer app.Close(appImpl)

	// Load the implicit Appfile
	implicit, err := appImpl.Implicit(appCtx)
	if err != nil {
		return nil, err
	}

	var final appfile.File
	if err := final.Merge(appDef); err != nil {
		return nil, fmt.Errorf("Error loading Appfile: %s", err)
	}
	if implicit != nil {
		if err := final.Merge(implicit); err != nil {
			return nil, fmt.Errorf("Error loading Appfile: %s", err)
		}
	}
	if f != nil {
		if err := final.Merge(f); err != nil {
			return nil, fmt.Errorf("Error loading Appfile: %s", err)
		}
	}

	return &final, nil
}
