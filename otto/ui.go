package otto

import (
	"io"

	"github.com/hashicorp/otto/ui"
)

// readerToUI takes an io.Reader and sends the data to a UI.
func readerToUI(uiVal ui.Ui, r io.Reader, doneCh chan<- struct{}) {
	defer close(doneCh)
	var buf [1024]byte
	for {
		n, err := r.Read(buf[:])
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
}
