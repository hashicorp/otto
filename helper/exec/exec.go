package exec

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/hashicorp/otto/ui"
)

// Run runs the given command and streams all the output to the
// given UI. It also connects stdin properly so that input works as
// expected.
func Run(uiVal ui.Ui, cmd *exec.Cmd) error {
	out_r, out_w := io.Pipe()
	cmd.Stdin = os.Stdin
	cmd.Stdout = out_w
	cmd.Stderr = out_w

	// Copy output to the UI until we can't.
	uiDone := make(chan struct{})
	go func() {
		defer close(uiDone)
		var buf [1024]byte
		for {
			n, err := out_r.Read(buf[:])
			if n > 0 {
				uiVal.Raw(string(buf[:n]))
			}

			// We just break on any error. io.EOF is not an error and
			// is our true exit case, but any other error we don't really
			// handle here. It probably means something went wrong
			// somewhere else anyways.
			if err != nil {
				break
			}
		}
	}()

	// Run the command
	log.Printf("[DEBUG] execDir: %s", cmd.Dir)
	log.Printf("[DEBUG] exec: %s %s", cmd.Path, strings.Join(cmd.Args[1:], " "))

	// Build a runnable command that we can log out to make things easier
	// for debugging. This lets debuging devs just copy and paste the command.
	var debugBuf bytes.Buffer
	for _, env := range cmd.Env {
		parts := strings.SplitN(env, "=", 2)
		debugBuf.WriteString(fmt.Sprintf("%s=%q ", parts[0], parts[1]))
	}
	debugBuf.WriteString(cmd.Path + " ")
	for _, arg := range cmd.Args[1:] {
		if strings.Contains(arg, " ") {
			debugBuf.WriteString(fmt.Sprintf("'%s' ", arg))
		} else {
			debugBuf.WriteString(fmt.Sprintf("%s ", arg))
		}
	}
	log.Printf("[DEBUG] exec runnable: %s", debugBuf.String())
	debugBuf.Reset()

	// Run
	err := cmd.Run()

	// Wait for all the output to finish
	out_w.Close()
	<-uiDone

	// Output one extra newline to separate output from Otto
	uiVal.Message("")

	// Return the output from the command
	return err
}

// OttoSkipCleanupEnvVar, when set, tells Otto to avoid cleaning up its
// temporary workspace files, which can be helpful for debugging.
const OttoSkipCleanupEnvVar = "OTTO_SKIP_CLEANUP"

// ShouldCleanup returns true for normal operation. It returns false if the
// user requested that Otto avoid cleaning up its temporary files for
// debugging purposes.
func ShouldCleanup() bool {
	return os.Getenv(OttoSkipCleanupEnvVar) == ""
}
