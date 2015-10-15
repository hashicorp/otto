package appfile

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/otto/appfile/detect"
	"github.com/hashicorp/terraform/dag"
)

var testHasGit bool

func init() {
	if _, err := exec.LookPath("git"); err == nil {
		testHasGit = true
	}
}

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

		{
			"compile-deps-detect",
			testCompileDepsStr,
			false,
		},

		{
			"compile-multi-dep",
			testCompileMultiDepStr,
			false,
		},

		{
			"compile-invalid",
			"",
			true,
		},

		{
			"compile-cycle",
			"",
			true,
		},

		{
			"compile-deps-no-id",
			"",
			true,
		},

		/*
			TODO: uncomment once we can enforce this
			{
				"compile-diff-project",
				"",
				true,
			},
		*/
	}

	for _, tc := range cases {
		t.Logf("Testing: %s", tc.Dir)

		// We wrap this in a function just so we can use defers
		func() {
			opts := testCompileOpts(t)
			defer os.RemoveAll(opts.Dir)
			f := testFile(t, tc.Dir)
			defer f.resetID()

			c, err := Compile(f, opts)
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s\n\n%s", tc.Dir, err)
			}

			if err == nil {
				testCompileCompare(t, c, tc.String)
				testCompileMarshal(t, c, opts.Dir)
			}
		}()
	}
}

// This is a really important test case that verifies that ".ottoid"
// is not ignored from dependencies. We had this happen with 0.1
func TestCompile_dotOttoId(t *testing.T) {
	if !testHasGit {
		t.Log("git not found, skipping")
		t.Skip()
	}

	opts := testCompileOpts(t)
	defer os.RemoveAll(opts.Dir)
	f := testFile(t, "compile-deps-git")
	defer f.resetID()

	// Rename DOTgit to .git since Git doesn't allow nested .git
	dir := filepath.Join(filepath.Dir(f.Path), "child")
	oldName := filepath.Join(dir, "DOTgit")
	newName := filepath.Join(dir, ".git")
	if err := os.Rename(oldName, newName); err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.Rename(newName, oldName)

	c, err := Compile(f, opts)
	if err != nil {
		t.Fatalf("err:\n\n%s", err)
	}

	testCompileCompare(t, c, testCompileDepGitStr)
	testCompileMarshal(t, c, opts.Dir)
}

func TestCompile_structure(t *testing.T) {
	cases := []struct {
		Dir  string
		Name string
		File *File
		Err  bool
	}{
		{
			"import-basic",
			"",
			&File{
				Application: &Application{
					Name: "foo",
					Type: "bar",
				},
				Project: &Project{
					Name:           "foo",
					Infrastructure: "aws",
				},
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
						Type: "aws",
					},
				},
			},
			false,
		},

		{
			"import-nested",
			"",
			&File{
				Application: &Application{
					Name: "bar",
					Type: "bar",
				},
				Project: &Project{
					Name:           "foo",
					Infrastructure: "aws",
				},
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
						Type: "aws",
					},
				},
			},
			false,
		},

		{
			"import-cycle",
			"",
			nil,
			true,
		},

		{
			"import-dep",
			"child",
			&File{
				Application: &Application{
					Name: "child",
					Type: "bar",
				},
				Project: &Project{
					Name:           "bar",
					Infrastructure: "aws",
				},
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "aws",
						Type: "aws",
					},
				},
			},
			false,
		},

		{
			"compile-dep-infra",
			"child",
			&File{
				Application: &Application{
					Name: "child",
					Type: "bar",
				},
				Project: &Project{
					Name:           "bar",
					Infrastructure: "google",
				},
				Infrastructure: []*Infrastructure{
					&Infrastructure{
						Name: "google",
						Type: "google",
					},
				},
			},
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
			f.initID()
			f.loadID()
			defer f.resetID()

			if tc.File != nil {
				tc.File.ID = f.ID
				tc.File.Path = f.Path
				tc.File.Imports = f.Imports
			}

			c, err := Compile(f, opts)
			if (err != nil) != tc.Err {
				t.Fatalf("err: %s\n\n%s", tc.Dir, err)
			}
			if err == nil && c == nil {
				t.Fatalf("bad: compiled is nil\n\n%s", tc.Dir)
			}
			if err != nil {
				return
			}

			// Get the root file. If we specified a name, then find
			// the dependent application with that name since that is
			// what we're comparing against.
			actual := c.File
			if tc.Name != "" {
				actual = nil
				c.Graph.Walk(func(raw dag.Vertex) error {
					v := raw.(*CompiledGraphVertex)
					if v.File.Application == nil {
						return nil
					}
					if v.File.Application.Name == tc.Name {
						actual = v.File
					}
					return nil
				})

				if actual == nil {
					t.Fatalf("err: %s\n\n%s not found in graph", tc.Dir, tc.Name)
				}

				// For child files, we just clear these out.
				actual.ID = ""
				actual.Path = ""
				actual.Source = ""
				actual.Imports = nil
				tc.File.ID = actual.ID
				tc.File.Path = actual.Path
				tc.File.Imports = actual.Imports
				tc.File.Source = actual.Source
			}

			if !reflect.DeepEqual(actual, tc.File) {
				t.Fatalf("err: %s\n\n%#v\n\n%#v", tc.Dir, actual, tc.File)
			}
		}()
	}
}

func TestCompileID(t *testing.T) {
	opts := testCompileOpts(t)
	defer os.RemoveAll(opts.Dir)
	f := testFile(t, "compile-id")
	defer f.resetID()

	c, err := Compile(f, opts)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if f.ID == "" {
		t.Fatalf("ID should not be blank")
	}
	if f.ID != c.File.ID {
		t.Fatalf("%s != %s", f.ID, c.File.ID)
	}
}

func TestCompileID_existing(t *testing.T) {
	opts := testCompileOpts(t)
	defer os.RemoveAll(opts.Dir)
	f := testFile(t, "compile-id-exists")
	if f.ID == "" {
		t.Fatalf("ID should not be blank")
	}

	copyId := f.ID
	c, err := Compile(f, opts)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if copyId != c.File.ID {
		t.Fatalf("%s != %s", f.ID, c.File.ID)
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

	return &CompileOpts{
		Dir:    dir,
		Detect: &detect.Config{},
	}
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

const testCompileDepGitStr = `
Compiled Appfile: %s

Dep Graph:
bar
foo
  bar
`

const testCompileMultiDepStr = `
Compiled Appfile: %s

Dep Graph:
bar
baz
foo
  bar
  baz
`
