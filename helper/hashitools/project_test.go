package hashitools

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/hashicorp/go-version"
)

func TestProjectLatestVersion(t *testing.T) {
	p := &Project{Name: "vagrant"}

	vsn, err := p.LatestVersion()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	min := version.Must(version.NewVersion("1.0.0"))
	max := version.Must(version.NewVersion("2.0.0"))
	if vsn.LessThan(min) || vsn.GreaterThan(max) {
		t.Fatalf("bad: %s", vsn)
	}
}

type stubInstaller struct {
	path string
}

func (t *stubInstaller) InstallAsk(installed, required, latest *version.Version) (bool, error) {
	return false, nil
}
func (t *stubInstaller) Install(*version.Version) error { return nil }
func (t *stubInstaller) Path() string                   { return t.path }

// https://github.com/hashicorp/otto/issues/70
func TestVersion_vagrantStdErrWarning(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Test uses a bash script; skipping on Windows")
	}

	path, err := filepath.Abs(filepath.Join(
		"./test-fixtures", "vagrant-version-stderr", "vagrant"))
	if err != nil {
		t.Fatal(err)
	}

	p := &Project{
		Name: "vagrant",
		Installer: &stubInstaller{
			path: path,
		},
	}

	v, err := p.Version()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if v == nil || v.String() != "1.2.3" {
		t.Fatalf("expected: 1.2.3, got: %s", v)
	}
}

func TestVersionRe(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{
			"0.8.7",
			"0.8.7",
		},

		{
			"Terraform v0.6.4-dev (4d37704d532ae3effdfe2c5b6254bee0b94e8d8e+CHANGES)\n",
			"0.6.4-dev",
		},

		{
			"Consul v0.5.0\nConsul Protocol: 2 (Understands back to: 1)",
			"0.5.0",
		},

		{
			"Packer v0.5.0.dev\n",
			"0.5.0.dev",
		},

		{
			"Vagrant 1.7.4",
			"1.7.4",
		},
	}

	for _, tc := range cases {
		matches := versionRe.FindStringSubmatch(tc.Input)
		if len(matches) == 0 {
			t.Fatalf("bad: %s", tc.Input)
		}
		if matches[1] != tc.Output {
			t.Fatalf("bad: %s != %s\n\n%s", matches[1], tc.Output, tc.Input)
		}
	}
}
