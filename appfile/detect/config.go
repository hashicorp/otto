package detect

import (
	"fmt"
	"path/filepath"
)

// Config is the format of the configuration files
type Config struct {
	Detectors []*Detector
}

// Merge merges another config into this one. This will modify this
// Config object. Detectors in c2 are tried after detectors in this
// Config. Conflicts are ignored as lower priority detectors, meaning that
// if two detectors are for type "go", both will be tried.
func (c *Config) Merge(c2 *Config) error {
	c.Detectors = append(c.Detectors, c2.Detectors...)
	return nil
}

// Detector is something that detects a single type.
type Detector struct {
	Type string
	File []string
}

// Detect will return true if this detector matches within the given
// directory.
func (d *Detector) Detect(dir string) (bool, error) {
	for _, pattern := range d.File {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return false, err
		}
		if len(matches) > 0 {
			return true, nil
		}
	}

	return false, nil
}

func (d *Detector) GoString() string {
	return fmt.Sprintf("*%#v", *d)
}
