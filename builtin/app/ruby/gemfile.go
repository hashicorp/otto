package rubyapp

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var (
	// regexp for finding gems in a gemfile
	gemGemfileRegexp = `\s*gem\s+['"]%s['"]`

	// regexp for finding gems in a gemfile.lock
	gemGemfileLockRegexp = `\s*%s\s+\(`
)

// HasGem checks if the Ruby project in the given directory has the
// specified gem. This uses Gemfile and Gemfile.lock to find this gem.
//
// If no Gemfile is in the directory, false is always returned.
func HasGem(dir, name string) (bool, error) {
	// Check normal Gemfile
	ok, err := hasGemGemfile(dir, name)
	if ok || err != nil {
		return ok, err
	}

	// Check Gemfile.lock
	ok, err = hasGemGemfileLock(dir, name)
	if ok || err != nil {
		return ok, err
	}

	// Nope!
	return false, nil
}

func hasGemGemfile(dir, name string) (bool, error) {
	path := filepath.Join(dir, "Gemfile")
	reStr := fmt.Sprintf(gemGemfileRegexp, name)
	return hasGemRaw(path, reStr)
}

func hasGemGemfileLock(dir, name string) (bool, error) {
	path := filepath.Join(dir, "Gemfile.lock")
	reStr := fmt.Sprintf(gemGemfileLockRegexp, name)
	return hasGemRaw(path, reStr)
}

func hasGemRaw(path, reStr string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}

		return false, err
	}

	// Try to find it!
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	re := regexp.MustCompile(reStr)
	return re.MatchReader(bufio.NewReader(f)), nil
}
