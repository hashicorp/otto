package rubyapp

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var rubyVersionGemfileRegexp = regexp.MustCompile(`ruby\s+['"]([.\d]+)['"]`)

// detectRubyVersion attempts to detect the Ruby version that needs to
// be installed by inspecting the environment (Gemfile, .ruby-version, etc.).
func detectRubyVersion(dir string) (result string, err error) {
	// Gemfile
	result, err = detectRubyVersionGemfile(dir)
	if result != "" || err != nil {
		return
	}

	// No version detected
	return "", nil
}

func detectRubyVersionGemfile(dir string) (result string, err error) {
	path := filepath.Join(dir, "Gemfile")

	// Verify the Gemfile exists
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return
	}

	// Open our file. We wrap the reader in a bufio so we can do a
	// streaming regexp
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	// Try to find a direct match.
	idx := rubyVersionGemfileRegexp.FindReaderSubmatchIndex(bufio.NewReader(f))
	if idx == nil {
		return
	}

	// Seek to the left of the first submatch
	if _, err = f.Seek(int64(idx[2]), 0); err != nil {
		return
	}

	resultBytes := make([]byte, idx[3]-idx[2])
	n, err := f.Read(resultBytes)
	if err != nil {
		return
	}
	if n != len(resultBytes) {
		err = fmt.Errorf("failed to read proper amount of bytes for Ruby version")
		return
	}

	result = string(resultBytes)
	return
}
