package appfile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompiled_impl(t *testing.T) {
	var _ json.Marshaler = new(Compiled)
	var _ json.Unmarshaler = new(Compiled)
}

func TestCompile(t *testing.T) {
	cases := []struct {
		Dir    string
		String string
		Err    bool
	}{
		{
			"compile-basic",
			testCompileBasicStr,
			false,
		},

		{
			"compile-deps",
			testCompileDepsStr,
			false,
		},
	}

	for _, tc := range cases {
		t.Logf("Testing: %s", tc.Dir)

		// We wrap this in a function just so we can use defers
		func() {
			opts := testCompileOpts(t)
			defer os.RemoveAll(opts.Dir)
			f := testFile(t, tc.Dir)

			c, err := Compile(f, opts)
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s\n\n%s", tc.Dir, err)
			}

			testCompileCompare(t, c, tc.String)
			testCompileMarshal(t, c, opts.Dir)
		}()
	}
}

func testCompileCompare(t *testing.T, c *Compiled, expected string) {
	actual := strings.TrimSpace(c.String())
	expected = strings.TrimSpace(fmt.Sprintf(expected, c.File.Path))
	if actual != expected {
		t.Fatalf("bad:\n\n%s\n\n%s", actual, expected)
	}
}

func testCompileMarshal(t *testing.T, original *Compiled, dir string) {
	c, err := LoadCompiled(dir)
	if err != nil {
		t.Fatalf("err loading compiled: %s", err)
	}

	if c.String() != original.String() {
		t.Fatalf("bad:\n\n%s\n\n%s", c, original)
	}
}

func testCompileOpts(t *testing.T) *CompileOpts {
	dir, err := ioutil.TempDir("", "otto-")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return &CompileOpts{Dir: dir}
}

func testFile(t *testing.T, dir string) *File {
	path := filepath.Join("./test-fixtures", dir, "Appfile")
	f, err := ParseFile(path)
	if err != nil {
		t.Fatalf("err: %s\n\n%s", path, err)
	}

	return f
}

const testCompileBasicStr = `
Compiled Appfile: %s

Dep Graph:
foo
`

const testCompileDepsStr = `
Compiled Appfile: %s

Dep Graph:
bar
foo
  bar
`
