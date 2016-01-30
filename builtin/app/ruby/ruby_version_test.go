package rubyapp

import (
	"path/filepath"
	"testing"
)

func TestDetectRubyVersion_gemfile(t *testing.T) {
	vsn, err := detectRubyVersion(filepath.Join("./test-fixtures", "ruby-version-gemfile"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if vsn != "2.2.2" {
		t.Fatalf("bad: %s", vsn)
	}
}

func TestDetectRubyVersion_gemfileNoPatch(t *testing.T) {
	vsn, err := detectRubyVersion(filepath.Join("./test-fixtures", "ruby-version-gemfile-nopatch"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if vsn != "2.3" {
		t.Fatalf("bad: %s", vsn)
	}
}

func TestDetectRubyVersion_gemfileNoVersion(t *testing.T) {
	vsn, err := detectRubyVersion(filepath.Join("./test-fixtures", "ruby-version-gemfile-none"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if vsn != "" {
		t.Fatalf("bad: %s", vsn)
	}
}

func TestDetectRubyVersion_rubyVersionFile(t *testing.T) {
	vsn, err := detectRubyVersion(filepath.Join("./test-fixtures", "ruby-version-file"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if vsn != "2.2.2" {
		t.Fatalf("bad: %s", vsn)
	}
}

func TestDetectRubyVersion_rubyVersionFileEmpty(t *testing.T) {
	vsn, err := detectRubyVersion(filepath.Join("./test-fixtures", "ruby-version-file-empty"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if vsn != "" {
		t.Fatalf("bad: %s", vsn)
	}
}

func TestDetectRubyVersion_rubyVersionPlusGemfile(t *testing.T) {
	vsn, err := detectRubyVersion(filepath.Join("./test-fixtures", "ruby-version-file-gemfile"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if vsn != "2.2.2" {
		t.Fatalf("bad: %s", vsn)
	}
}

func TestDetectRubyVersion_rubyVersionPlusGemfileNoPatch(t *testing.T) {
	vsn, err := detectRubyVersion(filepath.Join("./test-fixtures", "ruby-version-file-gemfile-nopatch"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	if vsn != "2.3" {
		t.Fatalf("bad: %s", vsn)
	}
}
