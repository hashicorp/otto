package appfile

import (
	"path/filepath"
	"testing"
)

func TestFileActiveInfrastructure(t *testing.T) {
	cases := []struct {
		File   string
		Result string
	}{
		{
			"file-active-infra-basic.hcl",
			"aws",
		},
	}

	for _, tc := range cases {
		path := filepath.Join("./test-fixtures", tc.File)
		actual, err := ParseFile(path)
		if err != nil {
			t.Fatalf("file: %s\n\n%s", tc.File, err)
			continue
		}

		infra := actual.ActiveInfrastructure()
		if infra.Name != tc.Result {
			t.Fatalf("file: %s\n\n%s", tc.File, infra.Name)
		}
	}
}
