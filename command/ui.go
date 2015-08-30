package command

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	"github.com/hashicorp/otto/ui"
	"github.com/hashicorp/vault/helper/password"
	"github.com/mitchellh/cli"
)

var defaultInputReader io.Reader
var defaultInputWriter io.Writer

// NewUi returns a new otto Ui implementation for use around
// the given CLI Ui implementation.
func NewUi(raw cli.Ui) ui.Ui {
	return &ui.Styled{
		Ui: &cliUi{
			CliUi: raw,
		},
	}
}

// cliUi is a wrapper around a cli.Ui that implements the otto.Ui
// interface. It is unexported since the NewUi method should be used
// instead.
type cliUi struct {
	CliUi cli.Ui

	// Reader and Writer are used for Input
	Reader io.Reader
	Writer io.Writer

	interrupted bool
	l           sync.Mutex
}

func (u *cliUi) Header(msg string) {
	u.CliUi.Output(ui.Colorize(msg))
}

func (u *cliUi) Message(msg string) {
	u.CliUi.Output(ui.Colorize(msg))
}

func (u *cliUi) Raw(msg string) {
	fmt.Print(msg)
}

func (i *cliUi) Input(opts *ui.InputOpts) (string, error) {
	r := i.Reader
	w := i.Writer
	if r == nil {
		r = defaultInputReader
	}
	if w == nil {
		w = defaultInputWriter
	}
	if r == nil {
		r = os.Stdin
	}
	if w == nil {
		w = os.Stdout
	}

	// Make sure we only ask for input once at a time. Terraform
	// should enforce this, but it doesn't hurt to verify.
	i.l.Lock()
	defer i.l.Unlock()

	// If we're interrupted, then don't ask for input
	if i.interrupted {
		return "", errors.New("interrupted")
	}

	// Listen for interrupts so we can cancel the input ask
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	// Build the output format for asking
	var buf bytes.Buffer
	buf.WriteString("[reset]")
	buf.WriteString(fmt.Sprintf("[bold]%s[reset]\n", opts.Query))
	if opts.Description != "" {
		s := bufio.NewScanner(strings.NewReader(opts.Description))
		for s.Scan() {
			buf.WriteString(fmt.Sprintf("  %s\n", s.Text()))
		}
		buf.WriteString("\n")
	}
	if opts.Default != "" {
		buf.WriteString("  [bold]Default:[reset] ")
		buf.WriteString(opts.Default)
		buf.WriteString("\n")
	}
	buf.WriteString("  [bold]Enter a value:[reset] ")

	// Ask the user for their input
	if _, err := fmt.Fprint(w, ui.Colorize(buf.String())); err != nil {
		return "", err
	}

	// Listen for the input in a goroutine. This will allow us to
	// interrupt this if we are interrupted (SIGINT)
	result := make(chan string, 1)
	if opts.Hide {
		f, ok := r.(*os.File)
		if !ok {
			return "", fmt.Errorf("reader must be a file")
		}

		line, err := password.Read(f)
		if err != nil {
			return "", err
		}

		result <- line
	} else {
		go func() {
			var line string
			if _, err := fmt.Fscanln(r, &line); err != nil {
				log.Printf("[ERR] UIInput scan err: %s", err)
			}

			result <- line
		}()
	}

	select {
	case line := <-result:
		fmt.Fprint(w, "\n")

		if line == "" {
			line = opts.Default
		}

		return line, nil
	case <-sigCh:
		// Print a newline so that any further output starts properly
		// on a new line.
		fmt.Fprintln(w)

		// Mark that we were interrupted so future Ask calls fail.
		i.interrupted = true

		return "", errors.New("interrupted")
	}
}
