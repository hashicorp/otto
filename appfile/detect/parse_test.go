package detect

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		File   string
		Result *Config
		Err    bool
	}{
		{
			"basic.hcl",
			&Config{
				Detectors: []*Detector{
					&Detector{
						Type: "go",
						File: []string{"*.go"},
					},
				},
			},
			false,
		},
	}

	for _, tc := range cases {
		path, err := filepath.Abs(filepath.Join("./test-fixtures", "parse", tc.File))
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tc.File, err)
			continue
		}

		actual, err := ParseFile(path)
		if (err != nil) != tc.Err {
			t.Fatalf("file: %s\n\n%s", tc.File, err)
			continue
		}

		if !reflect.DeepEqual(actual, tc.Result) {
			t.Fatalf("file: %s\n\n%#v\n\n%#v", tc.File, actual, tc.Result)
		}
	}
}

func TestParseDir(t *testing.T) {
	cases := []struct {
		Dir    string
		Result *Config
		Err    bool
	}{
		{
			"basic",
			&Config{
				Detectors: []*Detector{
					&Detector{
						Type: "go",
						File: []string{"*.go"},
					},
					&Detector{
						Type: "ruby",
						File: []string{"*.rb"},
					},
				},
			},
			false,
		},

		{
			"no-exist",
			nil,
			false,
		},
	}

	for _, tc := range cases {
		path, err := filepath.Abs(filepath.Join("./test-fixtures", "parse-dir", tc.Dir))
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tc.Dir, err)
			continue
		}

		actual, err := ParseDir(path)
		if (err != nil) != tc.Err {
			t.Fatalf("file: %s\n\n%s", tc.Dir, err)
			continue
		}

		if !reflect.DeepEqual(actual, tc.Result) {
			t.Fatalf("file: %s\n\n%#v\n\n%#v", tc.Dir, actual, tc.Result)
		}
	}
}
