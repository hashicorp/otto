package nodeapp

import (
	"path/filepath"
	"testing"
)

func TestDetectNodeVersion_jsonfile(t *testing.T) {
	vsn, err := detectNodeVersion(filepath.Join("./test-fixtures", "node-version-jsonfile"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if vsn != "0.10.3" {
		t.Fatalf("bad: %s", vsn)
	}
}

func TestDetectNodeVersion_jsonfileNoVersion(t *testing.T) {
	vsn, err := detectNodeVersion(filepath.Join("./test-fixtures", "node-version-jsonfile-none"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if vsn != "" {
		t.Fatalf("bad: %s", vsn)
	}
}
