package ui

import (
	"bufio"
	"bytes"
	"strings"
)

// Styled is a wrapper around an existing UI that automatically
// adds formatting around the UI text.
type Styled struct {
	Ui
}

func (u *Styled) Header(msg string) {
	u.Ui.Header(u.prefix("==> ", msg))
}

func (u *Styled) Message(msg string) {
	u.Ui.Message(u.prefix("    ", msg))
}

func (u *Styled) prefix(prefix, msg string) string {
	var buf bytes.Buffer

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
