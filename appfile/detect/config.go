package detect

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
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
	// Type is the type that will match if this detector matches
	Type string

	// File is a list of file globs to look for. If any are found, it is
	// a match.
	File []string

	// Contents is a content matcher. The key is a filename and the
	// path is the file contents regular expression.
	Contents map[string]string
}

// Detect will return true if this detector matches within the given
// directory.
func (d *Detector) Detect(dir string) (bool, error) {
	// First test files
	for _, pattern := range d.File {
		matches, err := filepath.Glob(filepath.Join(dir, pattern))
		if err != nil {
			return false, err
		}
		if len(matches) > 0 {
			return true, nil
		}
	}

	// Test contents
	for k, v := range d.Contents {
		path := filepath.Join(dir, k)
		if _, err := os.Stat(path); err != nil {
			continue
		}

		if ok, err := d.matchContents(path, v); err != nil {
			return false, err
		} else if ok {
			return true, nil
		}
	}

	return false, nil
}

func (d *Detector) matchContents(path string, raw string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	re, err := regexp.Compile(raw)
	if err != nil {
		return false, err
	}

	return re.MatchReader(bufio.NewReader(f)), nil
}

func (d *Detector) GoString() string {
	return fmt.Sprintf("*%#v", *d)
}
