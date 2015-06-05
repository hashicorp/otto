package directory

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// TestBackend is a public test helper that verifies a backend
// functions properly.
func TestBackend(t *testing.T, b Backend) {
	var buf bytes.Buffer
	var err error

	// PutBlob
	err = b.PutBlob("foo", &BlobData{Data: strings.NewReader("bar")})
	if err != nil {
		t.Fatalf("PutBlob error: %s", err)
	}

	// GetBlob (exists)
	data, err := b.GetBlob("foo")
	if err != nil {
		t.Fatalf("GetBlob error: %s", err)
	}
	_, err = io.Copy(&buf, data.Data)
	data.Close()
	if err != nil {
		t.Fatalf("GetBlob error: %s", err)
	}
	if buf.String() != "bar" {
		t.Fatalf("GetBlob bad data: %s", buf.String())
	}
}
