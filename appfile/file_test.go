package appfile

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestFileActiveInfrastructure(t *testing.T) {
	cases := []struct {
		File   string
		Result string
	}{
		{
			"file-active-infra-basic.hcl",
			"aws",
		},
	}

	for _, tc := range cases {
		path := filepath.Join("./test-fixtures", tc.File)
		actual, err := ParseFile(path)
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tc.File, err)
			continue
		}

		infra := actual.ActiveInfrastructure()
		if infra.Name != tc.Result {
			t.Fatalf("file: %s\n\n%s", tc.File, infra.Name)
		}
	}
}

func TestFileMerge(t *testing.T) {
	cases := map[string]struct {
		One, Two, Three *File
	}{
		"ID": {
			One: &File{
				ID: "foo",
			},
			Two: &File{
				ID: "bar",
			},
			Three: &File{
				ID: "bar",
			},
		},

		"Path": {
			One: &File{
				Path: "foo",
			},
			Two: &File{
				Path: "bar",
			},
			Three: &File{
				Path: "bar",
			},
		},

		"Application": {
			One: &File{
				Application: &Application{
					Name: "foo",
				},
			},
			Two: &File{
				Application: &Application{
					Type: "foo",
				},
			},
			Three: &File{
				Application: &Application{
					Name: "foo",
					Type: "foo",
				},
			},
		},

		"Application (no merge)": {
			One: &File{
				Application: &Application{
					Name: "foo",
				},
			},
			Two: &File{},
			Three: &File{
				Application: &Application{
					Name: "foo",
				},
			},
		},

		"Infra (no merge)": {
			One: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
					},
				},
			},
			Two: &File{},
			Three: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
					},
				},
			},
		},

		"Infra (add)": {
			One: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
					},
				},
			},
			Two: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "google",
					},
				},
			},
			Three: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
					},
					&Infrastructure{
						Name: "google",
					},
				},
			},
		},

		"Infra (override)": {
			One: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
					},
				},
			},
			Two: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
					},
				},
			},
			Three: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
					},
				},
			},
		},

		"Foundations (none)": {
			One: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
						Foundations: []*Foundation{
							&Foundation{
								Name: "consul",
							},
						},
					},
				},
			},
			Two: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
					},
				},
			},
			Three: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
						Foundations: []*Foundation{
							&Foundation{
								Name: "consul",
							},
						},
					},
				},
			},
		},

		"Foundations (override)": {
			One: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
						Foundations: []*Foundation{
							&Foundation{
								Name: "consul",
							},
						},
					},
				},
			},
			Two: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
						Foundations: []*Foundation{
							&Foundation{
								Name: "tubes",
							},
						},
					},
				},
			},
			Three: &File{
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
						Foundations: []*Foundation{
							&Foundation{
								Name: "tubes",
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		if err := tc.One.Merge(tc.Two); err != nil {
			t.Fatalf("%s: %s", name, err)
		}

		if !reflect.DeepEqual(tc.One, tc.Three) {
			t.Fatalf("%s:\n\n%#v\n\n%#v", name, tc.One, tc.Three)
		}
	}
}
