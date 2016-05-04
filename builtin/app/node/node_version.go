package nodeapp

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var nodeVersionJsonfileRegexp = regexp.MustCompile(`"node"\s:\s+['"]([.\d]+)['"]`)

// detectNodeVersion attempts to detect the Node version that needs to
// be installed by inspecting the package.json.
func detectNodeVersion(dir string) (result string, err error) {
	log.Printf("[DEBUG] node: Attempting to detect Node version for: %s", dir)

	// package.json
	result, err = detectNodeVersionJsonfile(dir)
	if result != "" || err != nil {
		return
	}

	// No version detected
	return "", nil
}

func detectNodeVersionJsonfile(dir string) (result string, err error) {
	path := filepath.Join(dir, "package.json")

	// Verify the package.json exists
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			log.Printf("[DEBUG] node: Nodefile not found, will not detect Node version this way")
			err = nil
		}

		return
	}

	log.Printf("[DEBUG] node: package.json found! Attempting to detect node version")

	// Open our file. We wrap the reader in a bufio so we can do a
	// streaming regexp
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	// Try to find a direct match.
	idx := nodeVersionJsonfileRegexp.FindReaderSubmatchIndex(bufio.NewReader(f))
	if idx == nil {
		log.Printf("[DEBUG] node: package.json has no 'node' declaration, no version found")
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
		err = fmt.Errorf("failed to read proper amount of bytes for node version")
		return
	}

	result = string(resultBytes)
	log.Printf("[DEBUG] node: package.json detected Node: %q", result)
	return
}
