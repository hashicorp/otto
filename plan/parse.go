package plan

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/mapstructure"
)

// ParseFile parses an HCL or JSON file into a set of plans.
func ParseFile(path string) ([]*Plan, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return Parse(f)
}

// Parse parses HCL or JSON data into a set of plans.
func Parse(r io.Reader) ([]*Plan, error) {
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
		"plan",
	}
	if err := checkHCLKeys(list, valid); err != nil {
		return nil, err
	}

	// Parse the plans
	if o := list.Filter("plan"); len(o.Items) > 0 {
		return parsePlans(o)
	}

	return nil, nil
}

func parsePlans(list *ast.ObjectList) ([]*Plan, error) {
	result := make([]*Plan, 0, len(list.Items))
	for i, item := range list.Items {
		// Check for invalid keys
		valid := []string{"description", "inputs", "task"}
		if err := checkHCLKeys(item.Val, valid); err != nil {
			return nil, multierror.Prefix(err, fmt.Sprintf("plan %d:", i+1))
		}

		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return nil, err
		}

		var plan Plan
		if err := mapstructure.WeakDecode(m, &plan); err != nil {
			return nil, err
		}

		// If we can parse tasks, parse those
		if ot, ok := item.Val.(*ast.ObjectType); ok {
			list := ot.List
			if o := list.Filter("task"); len(o.Items) > 0 {
				tasks, err := parseTasks(o)
				if err != nil {
					return nil, err
				}

				plan.Tasks = tasks
			}
		}

		result = append(result, &plan)
	}

	return result, nil
}

func parseTasks(list *ast.ObjectList) ([]*Task, error) {
	result := make([]*Task, 0, len(list.Items))
	for i, item := range list.Items {
		var t Task

		// Verify we have a key
		if len(item.Keys) != 1 {
			return nil, fmt.Errorf("task %d: needs exactly 1 key", i+1)
		}

		// Set our type
		t.Type = item.Keys[0].Token.Value().(string)

		// Decode into a map for ease
		var m map[string]interface{}
		if err := hcl.DecodeObject(&m, item.Val); err != nil {
			return nil, err
		}

		err := mapstructure.WeakDecode(m["description"], &t.Description)
		if err != nil {
			return nil, err
		}

		// Delete any keys that can't be args
		delete(m, "description")

		// Remainder are args
		if len(m) > 0 {
			t.Args = make(map[string]*TaskArg)
			for k, v := range m {
				t.Args[k] = &TaskArg{Value: v}
			}
		}

		result = append(result, &t)
	}

	return result, nil
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
