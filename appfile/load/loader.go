package load

import (
	"fmt"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/appfile/detect"
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

	// TODO: plugin loading and implicit
	return realFile, nil
}
