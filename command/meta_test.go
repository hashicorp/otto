package command

import (
	"strings"
	"testing"
)

func TestMetaRootDir(t *testing.T) {
	cases := []struct {
		Dir      string
		Err      bool
		Expected string
	}{
		{
			"rootdir-appfile",
			true,
			"",
		},

		{
			"rootdir-dototto",
			false,
			"",
		},
	}

	for _, tc := range cases {
		meta := TestMeta(t)

		dir := fixtureDir(tc.Dir)
		actual, err := meta.RootDir(dir)
		if (err != nil) != tc.Err {
			t.Fatalf("%s: err: %s", tc.Dir, err)
		}
		if tc.Err {
			continue
		}

		actual = strings.TrimPrefix(actual, dir)
		if actual != tc.Expected {
			t.Fatalf("%s: bad: %s", tc.Dir, actual)
		}
	}
}
