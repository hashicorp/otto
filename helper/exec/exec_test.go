package exec

import (
	"bytes"
	"os/exec"
	"runtime"
	"testing"

	"github.com/hashicorp/otto/ui"
)

func TestRun(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Windows. Not running this test.")
	}

	if _, err := exec.LookPath("echo"); err != nil {
		t.Skipf("echo not found, skipping test: %s", err)
	}

	cmd := exec.Command("echo", "-n", "hello world")
	ui := new(ui.Mock)

	err := Run(ui, cmd)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var output bytes.Buffer
	for _, v := range ui.RawBuf {
		output.WriteString(v)
	}

	if output.String() != "hello world" {
		t.Fatalf("bad: %s", output.String())
	}
}
