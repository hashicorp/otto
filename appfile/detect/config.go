package detect

import (
	"path/filepath"
)

// Config is the format of the configuration files
type Config struct {
	Detectors []*Detector
}

// Detector is something that detects a single type.
type Detector struct {
	Type string
	File string
}

// Detect will return true if this detector matches within the given
// directory.
func (d *Detector) Detect(dir string) (bool, error) {
	matches, err := filepath.Glob(filepath.Join(dir, d.File))
	if err != nil {
		return false, err
	}

	return len(matches) > 0, nil
}
