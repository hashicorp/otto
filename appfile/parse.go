package appfile

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/hcl"
	hclobj "github.com/hashicorp/hcl/hcl"
	"github.com/mitchellh/mapstructure"
)

// Parse parses the Appfile from the given io.Reader.
//
// Due to current internal limitations, the entire contents of the
// io.Reader will be copied into memory first before parsing.
func Parse(r io.Reader) (*File, error) {
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

	var result File

	// Parse the application
	if o := obj.Get("application", false); o != nil {
		if err := parseApplication(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'application': %s", err)
		}
	}

	// Parse the infrastructure
	if o := obj.Get("infrastructure", false); o != nil {
		if err := parseInfra(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'infrastructure': %s", err)
		}
	}

	return &result, nil
}

// ParseFile parses the given path as an Appfile.
func ParseFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

func parseApplication(result *File, obj *hclobj.Object) error {
	if obj.Len() > 1 {
		return fmt.Errorf("only one 'application' block allowed")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, obj); err != nil {
		return err
	}

	var app Application
	result.Application = &app
	return mapstructure.WeakDecode(m, &app)
}

func parseInfra(result *File, obj *hclobj.Object) error {
	// Get all the maps of keys to the actual object
	objects := make(map[string]*hclobj.Object)
	for _, o1 := range obj.Elem(false) {
		for _, o2 := range o1.Elem(true) {
			if _, ok := objects[o2.Key]; ok {
				return fmt.Errorf(
					"infrastructure '%s' defined more than once",
					o2.Key)
			}

			objects[o2.Key] = o2
		}
	}

	if len(objects) == 0 {
		return nil
	}

	// Go through each object and turn it into an actual result.
	collection := make([]*Infrastructure, 0, len(objects))
	for n, o := range objects {
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, o); err != nil {
			return err
		}

		var infra Infrastructure
		if err := mapstructure.WeakDecode(m, &infra); err != nil {
			return fmt.Errorf(
				"error parsing infrastructure '%s': %s", n, err)
		}

		infra.Name = n
		if infra.Type == "" {
			infra.Type = infra.Name
		}

		collection = append(collection, &infra)
	}

	result.Infrastructure = collection
	return nil
}
