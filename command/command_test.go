package command

import (
	"os"
	"path/filepath"
)

func init() {
	testingMode = true
}

func fixtureDir(n string) string {
	return filepath.Join("./test-fixtures", n)
}

func testEnv(k, v string) func() {
	old := os.Getenv(v)
	os.Setenv(k, v)
	return func() { os.Setenv(k, old) }
}
