package directory

import (
	"path/filepath"
	"reflect"
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
