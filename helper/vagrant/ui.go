package vagrant

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/armon/circbuf"
	"github.com/hashicorp/otto/ui"
)

// OutputCallback is the type that is called when there is a matching
// machine-readable output line for the UI.
type OutputCallback func(*Output)

// Output is a single line of machine-readable output from Vagrant.
type Output struct {
	Timestamp string
	Target    string
	Type      string
	Data      []string
}

// vagrantUi is an implementation of ui.Ui that we pass to helper/exec
// that parses the machine-readable output of Vagrant and calls callbacks
// for the various entries.
type vagrantUi struct {
	Callbacks map[string]OutputCallback

	buf  *circbuf.Buffer
	once sync.Once
}

// Finish should be called when you're done to force the final newline
// to parse the final event.
func (u *vagrantUi) Finish() {
	// We only need to do anything if we have stuff in our buffer
	if u.buf != nil && u.buf.TotalWritten() > 0 {
		u.Raw("\n")
	}
}

// Ignore these because we don't use them from helper/exec
func (u *vagrantUi) Header(string)                       {}
func (u *vagrantUi) Message(string)                      {}
func (u *vagrantUi) Input(*ui.InputOpts) (string, error) { return "", nil }

func (u *vagrantUi) Raw(msg string) {
	if msg == "" {
		// Not sure how this would happen, but there is nothing to do here.
		return
	}

	u.once.Do(u.init)

	// We loop while we have a newline in the message
	for {
		// Get the index for the newline if there is one
		idx := strings.IndexRune(msg, '\n')
		if idx == -1 {
			// The newline isn't there, write it to the circular buffer
			// and wait longer.
			u.buf.Write([]byte(msg))
			break
		}

		// We got a newline! Grab the contents from the circular buffer and
		// copy it so we can clear the buffer.
		bufRaw := u.buf.Bytes()
		buf := string(bufRaw) + msg[:idx]
		bufRaw = nil
		u.buf.Reset()

		// If we have more data, clear msg to that point. If we have
		// no more data, then just set message to empty and we'll break
		// on the next loop iteration
		if idx < len(msg) {
			msg = msg[idx+1:]
		} else {
			msg = ""
		}

		// Combine the data from the buffer up to the newline so we
		// have the full line, and split that by the commas.
		parts := strings.Split(buf, ",")
		if len(parts) < 3 {
			// Uh, invalid event?
			log.Printf("[ERROR] Invalid Vagrant event line: %s", buf)
			return
		}

		// Look for the callback
		cb, ok := u.Callbacks[parts[2]]
		if !ok {
			// No callback registered for this type, drop it
			return
		}

		// We have a callback, construct the output!
		var data []string
		if len(parts) > 3 {
			data = make([]string, len(parts)-3)
			for i, raw := range parts[3:] {
				data[i] = strings.Replace(
					strings.Replace(
						strings.Replace(raw, "%!(VAGRANT_COMMA)", ",", -1),
						"\\n", "\n", -1),
					"\\r", "\r", -1)
			}
		}

		// Callback
		cb(&Output{
			Timestamp: parts[0],
			Target:    parts[1],
			Type:      parts[2],
			Data:      data,
		})
	}
}

func (u *vagrantUi) init() {
	// Allocating the circular buffer. The circular buffer should only
	// keep track up to the point that there is a \n found so it doesn't
	// need to be huge, but it also limits the max length of an event.
	var err error
	u.buf, err = circbuf.NewBuffer(4096)
	if err != nil {
		panic(err)
	}
}

// For testing
func (o *Output) GoString() string {
	return fmt.Sprintf("*%#v", *o)
}
