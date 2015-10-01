package ui

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/mitchellh/colorstring"
)

// Styled is a wrapper around an existing UI that automatically
// adds formatting around the UI text.
type Styled struct {
	Ui
}

func (u *Styled) Header(msg string) {
	u.Ui.Header(u.prefix("[bold]==> ", msg))
}

func (u *Styled) Message(msg string) {
	u.Ui.Message(u.prefix("    ", msg))
}

func (u *Styled) prefix(prefix, msg string) string {
	var buf bytes.Buffer

	// We first write the color sequence (if any) of our message.
	// This makes it so that our prefix inherits the color property
	// of any message.
	buf.WriteString(colorstring.ColorPrefix(msg))

	scan := bufio.NewScanner(strings.NewReader(msg))
	for scan.Scan() {
		buf.WriteString(prefix)
		buf.WriteString(scan.Text() + "\n")
	}

	str := buf.String()
	if msg != "" {
		str = str[:len(str)-1]
	}

	return str
}
