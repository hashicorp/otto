package otto

import (
	"path/filepath"
)

func testPath(path ...string) string {
	args := make([]string, 1, len(path)+1)
	args[0] = "./test-fixtures"
	args = append(args, path...)

	return filepath.Join(args...)
}
