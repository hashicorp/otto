package load

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/appfile/detect"
)

func TestLoader_basic(t *testing.T) {
	cases := []struct {
		Path            string
		Input, Expected *appfile.File
	}{
		{
			"basic",
			&appfile.File{
				Application: &appfile.Application{
					Name: "foo",
				},
			},
			&appfile.File{
				Application: &appfile.Application{
					Name: "foo",
				},
			},
		},

		{
			"detect",
			&appfile.File{
				Application: &appfile.Application{
					Name: "foo",
				},
			},
			&appfile.File{
				Application: &appfile.Application{
					Name: "foo",
					Type: "test",
				},
			},
		},
	}

	l := testLoader(t)
	for _, tc := range cases {
		tc.Path = testPath(tc.Path)

		actual, err := l.Load(tc.Input, tc.Path)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		// Load the default and merge it
		def, err := appfile.Default(tc.Path, nil)
		if err != nil {
			t.Fatalf("err %s: %s", tc.Path, err)
		}
		if err := def.Merge(tc.Expected); err != nil {
			t.Fatalf("err %s: %s", tc.Path, err)
		}
		tc.Expected = def

		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("err: %s\n\n%#v", tc.Path, actual)
		}
	}
}

func testPath(path ...string) string {
	args := make([]string, len(path)+1)
	args[0] = "./test-fixtures"
	copy(args[1:], path)
	return filepath.Join(args...)
}

func testLoader(t *testing.T) *Loader {
	return &Loader{
		Detector: &detect.Config{
			Detectors: []*detect.Detector{
				&detect.Detector{
					Type: "test",
					File: []string{"test-file"},
				},
			},
		},
	}
}
