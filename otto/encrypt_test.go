package otto

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestCrypt(t *testing.T) {
	cases := []struct {
		Password string
		Data     []byte
	}{
		{
			"foo",
			[]byte("bar"),
		},
	}

	for _, tc := range cases {
		// Create a temporary file. We only need the path
		f, err := ioutil.TempFile("", "otto")
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		f.Close()
		path := f.Name()
		os.Remove(path)

		// Encrypt and decrypt!
		if err := cryptWrite(path, tc.Password, tc.Data); err != nil {
			t.Fatalf("write err: %s\n\n%s", tc.Password, err)
		}

		actual, err := cryptRead(path, tc.Password)
		if err != nil {
			t.Fatalf("read err: %s\n\n%s", tc.Password, err)
		}

		if !reflect.DeepEqual(actual, tc.Data) {
			t.Fatalf("read err: %s\n\n%#v\n\n%#v", tc.Password, actual, tc.Data)
		}
	}
}
