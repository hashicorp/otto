package rubyapp

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestHasGem(t *testing.T) {
	cases := []struct {
		Dir      string
		Gem      string
		Expected bool
	}{
		{
			"hasgem-basic",
			"redis",
			true,
		},

		{
			"hasgem-basic",
			"newp",
			false,
		},

		{
			"hasgem-lock-basic",
			"redis",
			true,
		},

		{
			"hasgem-lock-basic",
			"newp",
			false,
		},

		{
			"hasgem-empty",
			"newp",
			false,
		},
	}

	for _, tc := range cases {
		errPrefix := fmt.Sprintf("In '%s', looking for '%s': ", tc.Dir, tc.Gem)
		path := filepath.Join("./test-fixtures", tc.Dir)
		ok, err := HasGem(path, tc.Gem)
		if err != nil {
			t.Fatalf("%s: %s", errPrefix, err)
		}
		if ok != tc.Expected {
			t.Fatalf("%s: got %v", errPrefix, ok)
		}
	}
}
