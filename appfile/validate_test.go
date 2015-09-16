package appfile

import (
	"path/filepath"
	"testing"
)

func TestFileValidate(t *testing.T) {
	cases := []struct {
		File string
		Err  bool
	}{
		{
			"validate-basic",
			false,
		},

		{
			"validate-no-app",
			true,
		},

		{
			"validate-no-project",
			true,
		},

		{
			"validate-no-infra",
			true,
		},

		{
			"validate-app-no-name",
			true,
		},
	}

	for _, tc := range cases {
		f, err := ParseFile(filepath.Join("./test-fixtures", tc.File, "Appfile"))
		if err != nil {
			t.Fatalf("file:%s\n\n%s", tc.File, err)
			continue
		}

		err = f.Validate()
		if (err != nil) != tc.Err {
			t.Fatalf("file: %s\n\n%s", tc.File, err)
			continue
		}
	}
}
