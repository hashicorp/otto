package oneline

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// Read will read only the first line out of a file at the given path,
// stripping any whitespace from either side.
func Read(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	result, err := bufio.NewReader(f).ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	return strings.TrimSpace(result), nil
}
