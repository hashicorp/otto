package rubyapp

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/otto/helper/oneline"
)

var rubyVersionGemfileRegexp = regexp.MustCompile(`ruby\s+['"]([.\d]+)['"]`)

// detectRubyVersion attempts to detect the Ruby version that needs to
// be installed by inspecting the environment (Gemfile, .ruby-version, etc.).
func detectRubyVersion(dir string) (result string, err error) {
	log.Printf("[DEBUG] ruby: Attempting to detect Ruby version for: %s", dir)

	// .ruby-version
	result, err = detectRubyVersionFile(dir)
	if result != "" || err != nil {
		return
	}

	// Gemfile
	result, err = detectRubyVersionGemfile(dir)
	if result != "" || err != nil {
		return
	}

	// No version detected
	return "", nil
}

func detectRubyVersionFile(dir string) (result string, err error) {
	path := filepath.Join(dir, ".ruby-version")

	// Verify the .ruby-version exists
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Printf("[DEBUG] ruby: .ruby-version not found, will not detect Ruby version this way")
			err = nil
		}

		return
	}

	log.Printf("[DEBUG] ruby: .ruby-version found! Attempting to detect Ruby version")

	// Read the first line of the file
	result, err = oneline.Read(path)
	log.Printf("[DEBUG] ruby: Gemfile detected Ruby: %q", result)
	return
}

func detectRubyVersionGemfile(dir string) (result string, err error) {
	path := filepath.Join(dir, "Gemfile")

	// Verify the Gemfile exists
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Printf("[DEBUG] ruby: Gemfile not found, will not detect Ruby version this way")
			err = nil
		}

		return
	}

	log.Printf("[DEBUG] ruby: Gemfile found! Attempting to detect Ruby version")

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
		log.Printf("[DEBUG] ruby: Gemfile has no 'ruby' declaration, no version found")
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
	log.Printf("[DEBUG] ruby: Gemfile detected Ruby: %q", result)
	return
}
