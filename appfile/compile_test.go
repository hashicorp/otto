package appfile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompile_basic(t *testing.T) {
	opts := testCompileOpts(t)
	defer os.RemoveAll(opts.Dir)
	f := testFile(t, "compile-basic")

	c, err := Compile(f, opts)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testCompileCompare(t, c, testCompileBasicStr)
}

func TestCompile_deps(t *testing.T) {
	opts := testCompileOpts(t)
	defer os.RemoveAll(opts.Dir)
	f := testFile(t, "compile-deps")

	c, err := Compile(f, opts)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	testCompileCompare(t, c, testCompileDepsStr)
}

func testCompileCompare(t *testing.T, c *Compiled, expected string) {
	actual := strings.TrimSpace(c.String())
	expected = strings.TrimSpace(fmt.Sprintf(expected, c.File.Path))
	if actual != expected {
		t.Fatalf("bad:\n\n%s\n\n%s", actual, expected)
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
