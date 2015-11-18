package otto

import (
	"io"
)

// maybeClose is a function that can easily be deferred or called at any
// point and will close the given value if it is an io.Closer.
func maybeClose(v interface{}) error {
	if c, ok := v.(io.Closer); ok {
		return c.Close()
	}

	return nil
}
