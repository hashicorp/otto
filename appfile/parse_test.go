package appfile

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		File   string
		Result *File
		Err    bool
	}{
		{
			"basic.hcl",
			&File{
				Application: &Application{
					Name: "foo",
				},
				Project: &Project{
					Name:           "foo",
					Infrastructure: "aws",
					Stack: &Stack{
						Name: "bar",
					},
				},
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name:   "aws",
						Type:   "aws",
						Flavor: "foo",
					},
				},
			},
			false,
		},

		// Applications
		{
			"multi-app.hcl",
			nil,
			true,
		},

		// Infrastructures
		{
			"infra-dup.hcl",
			nil,
			true,
		},

		// Stacks
		{
			"multi-stack.hcl",
			nil,
			true,
		},
	}

	for _, tc := range cases {
		path := filepath.Join("./test-fixtures", tc.File)
		actual, err := ParseFile(path)
		if (err != nil) != tc.Err {
			t.Fatalf("file: %s\n\n%s", tc.File, err)
			continue
		}

		if actual != nil {
			if actual.Path != path {
				t.Fatalf("file: %s\n\n%s", tc.File, actual.Path)
			}
			actual.Path = ""
		}

		if !reflect.DeepEqual(actual, tc.Result) {
			t.Fatalf("file: %s\n\n%#v\n\n%#v", tc.File, actual, tc.Result)
		}
	}
}
