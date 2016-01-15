package directory

import (
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/otto/appfile"
)

func TestNewAppCompiled(t *testing.T) {
	cases := []struct {
		Input    string
		Err      bool
		Expected *App
	}{
		{
			"app-basic",
			false,
			&App{
				AppLookup: AppLookup{
					AppID:   "65972869-4e36-6344-b2ca-cd34cba1d3f7",
					Version: "1.0.0",
				},
				Name: "foo",
				Type: "bar",
			},
		},

		{
			"app-deps",
			false,
			&App{
				AppLookup: AppLookup{
					AppID:   "65972869-4e36-6344-b2ca-cd34cba1d3f7",
					Version: "1.0.0",
				},
				Name: "foo",
				Type: "bar",

				Dependencies: []AppLookup{
					AppLookup{
						AppID:   "hello",
						Version: "0.1.0",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		path := filepath.Join("./test-fixtures", tc.Input, "Appfile")
		c := appfile.TestAppfile(t, path)
		root, err := c.Graph.Root()
		if err != nil {
			t.Fatalf("%s: err: %s", tc.Input, err)
		}

		actual, err := NewAppCompiled(c, root)
		if (err != nil) != tc.Err {
			t.Fatalf("%s: err: %s", tc.Input, err)
		}
		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("%s: bad: %#v", tc.Input, actual)
		}
	}
}

func TestAppSlice_impl(t *testing.T) {
	var _ sort.Interface = new(AppSlice)
}

func TestAppSlice_sort(t *testing.T) {
	cases := []struct {
		Desc     string
		Input    []*App
		Expected []*App
	}{
		{
			"basic, name only",
			[]*App{
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "foo",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "bar",
				},
			},

			[]*App{
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "bar",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "foo",
				},
			},
		},

		{
			"name, then ID",
			[]*App{
				&App{
					AppLookup: AppLookup{
						AppID:   "b",
						Version: "1.0.0",
					},

					Name: "foo",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "foo",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "bar",
				},
			},

			[]*App{
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "bar",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "foo",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "b",
						Version: "1.0.0",
					},

					Name: "foo",
				},
			},
		},

		{
			"name, then ID, then version",
			[]*App{
				&App{
					AppLookup: AppLookup{
						AppID:   "b",
						Version: "1.2.3",
					},

					Name: "foo",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "b",
						Version: "1.10.0",
					},

					Name: "foo",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "foo",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "bar",
				},
			},

			[]*App{
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "bar",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "a",
						Version: "1.0.0",
					},

					Name: "foo",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "b",
						Version: "1.2.3",
					},

					Name: "foo",
				},
				&App{
					AppLookup: AppLookup{
						AppID:   "b",
						Version: "1.10.0",
					},

					Name: "foo",
				},
			},
		},
	}

	for _, tc := range cases {
		sort.Sort(AppSlice(tc.Input))
		if !reflect.DeepEqual(tc.Expected, tc.Input) {
			t.Fatalf("%s bad: %#v", tc.Desc, tc.Input)
		}
	}
}
