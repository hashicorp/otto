package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DevDep has information about an upstream dependency that should be
// used by the Dev function in order to build a complete development
// environment.
type DevDep struct {
	// Files is a list of files that this dependency created or uses.
	// If these files already exist, then future DevDep calls won't be
	// called and the cached data will be used.
	//
	// All files in this must be in the Context.CacheDir. Relative paths
	// will be expanded relative to the CacheDir. If the file is not
	// in the CacheDir, no caching will occur. The log will note if this
	// is happening.
	Files []string `json:"files"`
}

// RelFiles makes all the Files values relative to the given directory.
func (d *DevDep) RelFiles(dir string) error {
	for i, f := range d.Files {
		// If the path is already relative, ignore it
		if !filepath.IsAbs(f) {
			continue
		}

		// Make the path relative
		f, err := filepath.Rel(dir, f)
		if err != nil {
			return fmt.Errorf(
				"couldn't make directory relative: %s\n\n%s",
				d.Files[i], err)
		}

		d.Files[i] = f
	}

	return nil
}

// ReadDevDep reads a marshalled DevDep from disk.
func ReadDevDep(path string) (*DevDep, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var result DevDep
	dec := json.NewDecoder(f)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// WriteDevDep writes a DevDep out to disk.
func WriteDevDep(path string, dep *DevDep) error {
	// Pretty-print the JSON data so that it can be more easily inspected
	data, err := json.MarshalIndent(dep, "", "    ")
	if err != nil {
		return err
	}

	// Write it out
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(data))
	return err
}
