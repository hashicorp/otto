package rpc

import (
	"io"
	"log"
)

// serveSingleCopy is a helper that creates a side-channel on our yamux
// connection to send a stream of raw data.
//
// It is very important to wait for this to complete by listening on the
// doneCh that is sent in so data isn't corrupted.
func serveSingleCopy(
	name string,
	mux *muxBroker,
	doneCh chan<- struct{},
	id uint32, dst io.Writer, src io.Reader) {
	defer close(doneCh)

	conn, err := mux.Accept(id)
	if err != nil {
		log.Printf("[ERR] '%s' accept error: %s", name, err)
		return
	}

	// Be sure to close the connection after we're done copying so
	// that an EOF will successfully be sent to the remote side
	defer conn.Close()

	// The connection is the destination/source that is nil
	if dst == nil {
		dst = conn
	} else {
		src = conn
	}

	written, err := io.Copy(dst, src)
	log.Printf("[INFO] %d bytes written for '%s'", written, name)
	if err != nil {
		log.Printf("[ERR] '%s' copy error: %s", name, err)
	}
}
