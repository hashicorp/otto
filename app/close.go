package app

import (
	"io"
)

// Close is a function that can easily be deferred or called at any
// point and will close the given value if it is an io.Closer.
//
// This should be called with App implementations since they may be
// happening over RPC.
func Close(v interface{}) error {
	if c, ok := v.(io.Closer); ok {
		return c.Close()
	}

	return nil
}
