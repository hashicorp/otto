package scriptpack

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	testBasic "github.com/hashicorp/otto/scriptpack/test-basic"
)

func TestScriptPackWrite(t *testing.T) {
	// Build the SP
	sp := &ScriptPack{
		Name: "foo",
		Data: testBasic.Bindata,
	}

	// Temporary dir to test contents
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	// Write the data!
	if err := sp.Write(td); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify some files exist
	var path string
	path = filepath.Join(td, sp.Name, "hello.txt")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestScriptPackWriteArchive(t *testing.T) {
	// Build the SP
	sp := &ScriptPack{
		Name: "foo",
		Data: testBasic.Bindata,
	}

	// Temporary dir to test contents
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	// Write the data!
	archivePath := filepath.Join(td, "foo.tar.gz")
	if err := sp.WriteArchive(archivePath); err != nil {
		t.Fatalf("err: %s", err)
	}

	// Verify some files exist
	files := testArchive(t, archivePath, false)
	expected := []string{"foo/", "foo/hello.txt", "main.sh"}
	if !reflect.DeepEqual(files, expected) {
		t.Fatalf("bad: %#v", files)
	}
}

func testArchive(t *testing.T, path string, detailed bool) []string {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer f.Close()

	gzipR, err := gzip.NewReader(f)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tarR := tar.NewReader(gzipR)

	// Read all the entries
	result := make([]string, 0, 5)
	for {
		hdr, err := tarR.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		text := hdr.Name
		if detailed {
			// Check if the file is executable. We use these stub names
			// to compensate for umask differences in test environments
			// and limitations in using "git clone".
			if hdr.FileInfo().Mode()&0111 != 0 {
				text = hdr.Name + "-exec"
			} else {
				text = hdr.Name + "-reg"
			}
		}

		result = append(result, text)
	}

	sort.Strings(result)
	return result
}
