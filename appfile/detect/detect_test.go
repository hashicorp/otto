package detect

import (
	"path/filepath"
	"testing"
)

func TestApp(t *testing.T) {
	cases := []struct {
		Dir       string
		Expected  string
		Err       bool
		Detectors []*Detector
	}{
		{
			"app-none",
			"",
			false,
			[]*Detector{
				&Detector{
					Type: "go",
					File: "*.go",
				},
			},
		},

		{
			"app-go",
			"go",
			false,
			[]*Detector{
				&Detector{
					Type: "go",
					File: "*.go",
				},
			},
		},
	}

	for _, tc := range cases {
		actual, err := App(filepath.Join("test-fixtures", tc.Dir), &Config{
			Detectors: tc.Detectors,
		})
		if (err != nil) != tc.Err {
			t.Fatalf("%s err: %s", tc.Dir, err)
		}

		if actual != tc.Expected {
			t.Fatalf("%s: %s != %s", tc.Dir, actual, tc.Expected)
		}
	}
}
