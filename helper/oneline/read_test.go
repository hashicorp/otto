package oneline

import (
	"path/filepath"
	"testing"
)

func TestRead(t *testing.T) {
	cases := []struct {
		Name     string
		Expected string
		Err      bool
	}{
		{
			"basic.txt",
			"foo",
			false,
		},

		{
			"oneline.txt",
			"bar",
			false,
		},
	}

	for _, tc := range cases {
		actual, err := Read(filepath.Join("test-fixtures", tc.Name))
		if (err != nil) != tc.Err {
			t.Fatalf("%s, err: %s", tc.Name, err)
		}

		if actual != tc.Expected {
			t.Fatalf("%s, mismatch: %s != %s", tc.Name, actual, tc.Expected)
		}
	}
}
