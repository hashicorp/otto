package appfile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
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

	// Check for invalid keys
	valid := []string{
		"application",
		"customization",
		"import",
		"infrastructure",
		"project",
	}
	if err := checkHCLKeys(obj, valid); err != nil {
		return nil, err
	}

	var result File

	// Parse the imports
	if o := obj.Get("import", false); o != nil {
		if err := parseImport(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'import': %s", err)
		}
	}

	// Parse the application
	if o := obj.Get("application", false); o != nil {
		if err := parseApplication(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'application': %s", err)
		}
	}

	// Parse the project
	if o := obj.Get("project", false); o != nil {
		if err := parseProject(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'project': %s", err)
		}
	}

	// Parse the infrastructure
	if o := obj.Get("infrastructure", false); o != nil {
		if err := parseInfra(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'infrastructure': %s", err)
		}
	}

	// Parse the customizations
	if o := obj.Get("customization", false); o != nil {
		if err := parseCustomizations(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'customization': %s", err)
		}
	}

	return &result, nil
}

// ParseFile parses the given path as an Appfile.
func ParseFile(path string) (*File, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result, err := Parse(f)
	if result != nil {
		result.Path = path
		if err := result.loadID(); err != nil {
			return nil, err
		}
	}

	return result, err
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

func parseCustomizations(result *File, obj *hclobj.Object) error {
	// Get all the maps of keys to the actual object
	objects := make(map[string]*hclobj.Object)
	for _, o1 := range obj.Elem(false) {
		for _, o2 := range o1.Elem(true) {
			if _, ok := objects[o2.Key]; ok {
				return fmt.Errorf(
					"customization '%s' defined more than once",
					o2.Key)
			}

			objects[o2.Key] = o2
		}
	}

	if len(objects) == 0 {
		return nil
	}

	// Go through each object and turn it into an actual result.
	collection := make([]*Customization, 0, len(objects))
	for n, o := range objects {
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, o); err != nil {
			return err
		}

		var c Customization
		c.Type = strings.ToLower(n)
		c.Config = m

		collection = append(collection, &c)
	}

	result.Customization = &CustomizationSet{Raw: collection}
	return nil
}

func parseImport(result *File, obj *hclobj.Object) error {
	// Get all the maps of keys to the actual object
	objects := make([]*hclobj.Object, 0, 3)
	set := make(map[string]struct{})
	for _, o1 := range obj.Elem(false) {
		for _, o2 := range o1.Elem(true) {
			if _, ok := set[o2.Key]; ok {
				return fmt.Errorf(
					"imported '%s' more than once",
					o2.Key)
			}

			objects = append(objects, o2)
			set[o2.Key] = struct{}{}
		}
	}

	if len(objects) == 0 {
		return nil
	}

	// Go through each object and turn it into an actual result.
	collection := make([]*Import, 0, len(objects))
	for _, o := range objects {
		collection = append(collection, &Import{
			Source: o.Key,
		})
	}

	result.Imports = collection
	return nil
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

		// Parse the foundations if we have any
		if o2 := o.Get("foundation", false); o != nil {
			if err := parseFoundations(&infra, o2); err != nil {
				return fmt.Errorf("error parsing 'foundation': %s", err)
			}
		}

		collection = append(collection, &infra)
	}

	result.Infrastructure = collection
	return nil
}

func parseFoundations(result *Infrastructure, obj *hclobj.Object) error {
	// Get all the maps of keys to the actual object
	objects := make(map[string]*hclobj.Object)
	for _, o1 := range obj.Elem(false) {
		for _, o2 := range o1.Elem(true) {
			if _, ok := objects[o2.Key]; ok {
				return fmt.Errorf(
					"foundation '%s' defined more than once",
					o2.Key)
			}

			objects[o2.Key] = o2
		}
	}

	if len(objects) == 0 {
		return nil
	}

	// Go through each object and turn it into an actual result.
	collection := make([]*Foundation, 0, len(objects))
	for n, o := range objects {
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, o); err != nil {
			return err
		}

		var f Foundation
		f.Name = n
		f.Config = m

		collection = append(collection, &f)
	}

	// Set the results
	result.Foundations = collection
	return nil
}

func parseProject(result *File, obj *hclobj.Object) error {
	if obj.Len() > 1 {
		return fmt.Errorf("only one 'project' block allowed")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, obj); err != nil {
		return err
	}

	// Parse the project
	var proj Project
	result.Project = &proj
	if err := mapstructure.WeakDecode(m, &proj); err != nil {
		return err
	}

	return nil
}

func checkHCLKeys(obj *hclobj.Object, valid []string) error {
	validMap := make(map[string]struct{}, len(valid))
	for _, v := range valid {
		validMap[v] = struct{}{}
	}

	var result error
	for _, o := range obj.Elem(true) {
		if _, ok := validMap[o.Key]; !ok {
			result = multierror.Append(result, fmt.Errorf(
				"invald key: %s", o.Key))
		}
	}

	return result
}
