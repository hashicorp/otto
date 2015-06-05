package command

import (
	"path/filepath"
)

func fixtureDir(n string) string {
	return filepath.Join("./test-fixtures", n)
}
