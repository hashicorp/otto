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
					Name:   "foo",
					Detect: true,
					Dependencies: []*Dependency{
						&Dependency{
							Source: "foo",
						},
						&Dependency{
							Source: "bar",
						},
					},
				},
				Project: &Project{
					Name:           "foo",
					Infrastructure: "aws",
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

		{
			"app-no-detect.hcl",
			&File{
				Application: &Application{
					Name:   "foo",
					Detect: false,
				},
			},
			false,
		},

		// Customizations
		{
			"basic-custom.hcl",
			&File{
				Customization: &CustomizationSet{
					Raw: []*Customization{
						&Customization{
							Type: "dev",
							Config: map[string]interface{}{
								"go_version": "1.5",
							},
						},
					},
				},
			},
			false,
		},

		{
			"basic-custom-no-name.hcl",
			&File{
				Customization: &CustomizationSet{
					Raw: []*Customization{
						&Customization{
							Type: "app",
							Config: map[string]interface{}{
								"go_version": "1.5",
							},
						},
					},
				},
			},
			false,
		},

		{
			"basic-custom-case.hcl",
			&File{
				Customization: &CustomizationSet{
					Raw: []*Customization{
						&Customization{
							Type: "dev",
							Config: map[string]interface{}{
								"go_version": "1.5",
							},
						},
					},
				},
			},
			false,
		},

		// Infrastructures
		{
			"infra-dup.hcl",
			nil,
			true,
		},

		{
			"infra-foundations.hcl",
			&File{
				Application: &Application{
					Name:   "foo",
					Detect: true,
				},
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name:   "aws",
						Type:   "aws",
						Flavor: "foo",
						Foundations: []*Foundation{
							&Foundation{
								Name: "consul",
								Config: map[string]interface{}{
									"foo": "bar",
								},
							},
						},
					},
				},
			},
			false,
		},

		{
			"infra-foundations-dup.hcl",
			nil,
			true,
		},

		// Imports

		{
			"imports.hcl",
			&File{
				Application: &Application{
					Name:   "otto",
					Type:   "go",
					Detect: true,
				},
				Imports: []*Import{
					&Import{
						Source: "./foo",
					},
				},
			},
			false,
		},

		// Unknown keys
		{
			"unknown-keys.hcl",
			nil,
			true,
		},
	}

	for _, tc := range cases {
		path, err := filepath.Abs(filepath.Join("./test-fixtures", tc.File))
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tc.File, err)
			continue
		}

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
