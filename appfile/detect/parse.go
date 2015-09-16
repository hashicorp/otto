package detect

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl"
	hclobj "github.com/hashicorp/hcl/hcl"
	"github.com/mitchellh/mapstructure"
)

// Parse parses the detector config from the given io.Reader.
//
// Due to current internal limitations, the entire contents of the
// io.Reader will be copied into memory first before parsing.
func Parse(r io.Reader) (*Config, error) {
	// Copy the reader into an in-memory buffer first since HCL requires it.
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}

	// Parse the buffer
	obj, err := hcl.Parse(buf.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing: %s", err)
	}
	buf.Reset()

	var result Config

	// Parse the detects
	if o := obj.Get("detect", false); o != nil {
		if err := parseDetect(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'import': %s", err)
		}
	}

	return &result, nil
}

// ParseFile parses the given path as a single detector config.
func ParseFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

func parseDetect(result *Config, obj *hclobj.Object) error {
	// Get all the maps of keys to the actual object
	objects := make([]*hclobj.Object, 0, 2)
	for _, o1 := range obj.Elem(false) {
		for _, o2 := range o1.Elem(true) {
			objects = append(objects, o2)
		}
	}

	if len(objects) == 0 {
		return nil
	}

	// Go through each object and turn it into an actual result.
	collection := make([]*Detector, 0, len(objects))
	for _, o := range objects {
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, o); err != nil {
			return err
		}

		var d Detector
		if err := mapstructure.WeakDecode(m, &d); err != nil {
			return fmt.Errorf(
				"error parsing detector '%s': %s", o.Key, err)
		}

		d.Type = o.Key
		collection = append(collection, &d)
	}

	result.Detectors = collection
	return nil
}
