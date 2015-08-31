package directory

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestBlobDataWriteToFile(t *testing.T) {
	tf, err := ioutil.TempFile("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	tf.Close()
	defer os.Remove(tf.Name())

	raw := "foo"
	data := &BlobData{
		Key:  "bar",
		Data: strings.NewReader(raw),
	}

	if err := data.WriteToFile(tf.Name()); err != nil {
		t.Fatalf("err: %s", err)
	}

	actual, err := ioutil.ReadFile(tf.Name())
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if string(actual) != raw {
		t.Fatalf("bad: %s", actual)
	}
}
