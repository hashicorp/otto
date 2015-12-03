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
	"github.com/hashicorp/hcl/hcl/ast"
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
	root, err := hcl.Parse(buf.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing: %s", err)
	}
	buf.Reset()

	// Top-level item should be the object list
	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, fmt.Errorf("error parsing: file doesn't contain a root object")
	}

	// Check for invalid keys
	valid := []string{
		"application",
		"customization",
		"import",
		"infrastructure",
		"project",
	}
	if err := checkHCLKeys(list, valid); err != nil {
		return nil, err
	}

	var result File

	// Parse the imports
	if o := list.Filter("import"); len(o.Items) > 0 {
		if err := parseImport(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'import': %s", err)
		}
	}

	// Parse the application
	if o := list.Filter("application"); len(o.Items) > 0 {
		if err := parseApplication(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'application': %s", err)
		}
	}

	// Parse the project
	if o := list.Filter("project"); len(o.Items) > 0 {
		if err := parseProject(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'project': %s", err)
		}
	}

	// Parse the infrastructure
	if o := list.Filter("infrastructure"); len(o.Items) > 0 {
		if err := parseInfra(&result, o); err != nil {
			return nil, fmt.Errorf("error parsing 'infrastructure': %s", err)
		}
	}

	// Parse the customizations
	if o := list.Filter("customization"); len(o.Items) > 0 {
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

func parseApplication(result *File, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'application' block allowed")
	}

	// Get our one item
	item := list.Items[0]

	// Check for invalid keys
	valid := []string{"name", "type", "detect", "dependency"}
	if err := checkHCLKeys(item.Val, valid); err != nil {
		return multierror.Prefix(err, "application:")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
		return err
	}

	app := Application{Detect: true}
	result.Application = &app
	return mapstructure.WeakDecode(m, &app)
}

func parseCustomizations(result *File, list *ast.ObjectList) error {
	// Go through each object and turn it into an actual result.
	collection := make([]*Customization, 0, len(list.Items))
	for _, item := range list.Items {
		var key string
		if len(item.Keys) > 0 {
			key = item.Keys[0].Token.Value().(string)
		}

		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return err
		}

		var c Customization
		c.Name = strings.ToLower(key)
		c.Config = m

		collection = append(collection, &c)
	}

	result.Customization = &CustomizationSet{Raw: collection}
	return nil
}

func parseImport(result *File, list *ast.ObjectList) error {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil
	}

	// Go through each object and turn it into an actual result.
	collection := make([]*Import, 0, len(list.Items))
	seen := make(map[string]struct{})
	for _, item := range list.Items {
		key := item.Keys[0].Token.Value().(string)

		// Make sure we haven't already found this import
		if _, ok := seen[key]; ok {
			return fmt.Errorf("import '%s' defined more than once", key)
		}
		seen[key] = struct{}{}

		// Check for invalid keys
		if err := checkHCLKeys(item.Val, nil); err != nil {
			return multierror.Prefix(err, fmt.Sprintf(
				"import '%s':", key))
		}

		collection = append(collection, &Import{
			Source: key,
		})
	}

	result.Imports = collection
	return nil
}

func parseInfra(result *File, list *ast.ObjectList) error {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil
	}

	// Go through each object and turn it into an actual result.
	collection := make([]*Infrastructure, 0, len(list.Items))
	seen := make(map[string]struct{})
	for _, item := range list.Items {
		n := item.Keys[0].Token.Value().(string)

		// Make sure we haven't already found this
		if _, ok := seen[n]; ok {
			return fmt.Errorf("infrastructure '%s' defined more than once", n)
		}
		seen[n] = struct{}{}

		// Check for invalid keys
		valid := []string{"name", "type", "flavor", "foundation"}
		if err := checkHCLKeys(item.Val, valid); err != nil {
			return multierror.Prefix(err, fmt.Sprintf(
				"infrastructure '%s':", n))
		}

		var listVal *ast.ObjectList
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			listVal = ot.List
		} else {
			return fmt.Errorf("infrastructure '%s': should be an object", n)
		}

		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
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
		if o2 := listVal.Filter("foundation"); len(o2.Items) > 0 {
			if err := parseFoundations(&infra, o2); err != nil {
				return fmt.Errorf("error parsing 'foundation': %s", err)
			}
		}

		collection = append(collection, &infra)
	}

	result.Infrastructure = collection
	return nil
}

func parseFoundations(result *Infrastructure, list *ast.ObjectList) error {
	list = list.Children()
	if len(list.Items) == 0 {
		return nil
	}

	// Go through each object and turn it into an actual result.
	collection := make([]*Foundation, 0, len(list.Items))
	seen := make(map[string]struct{})
	for _, item := range list.Items {
		n := item.Keys[0].Token.Value().(string)

		// Make sure we haven't already found this
		if _, ok := seen[n]; ok {
			return fmt.Errorf("foundation '%s' defined more than once", n)
		}
		seen[n] = struct{}{}

		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
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

func parseProject(result *File, list *ast.ObjectList) error {
	if len(list.Items) > 1 {
		return fmt.Errorf("only one 'project' block allowed")
	}

	// Get our one item
	item := list.Items[0]

	// Check for invalid keys
	valid := []string{"name", "infrastructure"}
	if err := checkHCLKeys(item.Val, valid); err != nil {
		return multierror.Prefix(err, "project:")
	}

	var m map[string]interface{}
	if err := hcl.DecodeObject(&m, item.Val); err != nil {
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

func checkHCLKeys(node ast.Node, valid []string) error {
	var list *ast.ObjectList
	switch n := node.(type) {
	case *ast.ObjectList:
		list = n
	case *ast.ObjectType:
		list = n.List
	default:
		return fmt.Errorf("cannot check HCL keys of type %T", n)
	}

	validMap := make(map[string]struct{}, len(valid))
	for _, v := range valid {
		validMap[v] = struct{}{}
	}

	var result error
	for _, item := range list.Items {
		key := item.Keys[0].Token.Value().(string)
		if _, ok := validMap[key]; !ok {
			result = multierror.Append(result, fmt.Errorf(
				"invalid key: %s", key))
		}
	}

	return result
}
