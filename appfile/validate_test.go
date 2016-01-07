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

		{
			"validate-app-no-type",
			true,
		},

		{
			"validate-app-version",
			false,
		},

		{
			"validate-app-bad-version",
			true,
		},

		{
			"validate-project-no-name",
			true,
		},

		{
			"validate-project-no-infra",
			true,
		},

		{
			"validate-project-unknown-infra",
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
