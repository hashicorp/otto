package command

import (
	"path/filepath"
)

func init() {
	testingMode = true
}

func fixtureDir(n string) string {
	return filepath.Join("./test-fixtures", n)
}
