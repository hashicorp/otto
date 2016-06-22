package terraform

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestOutputs(t *testing.T) {
	cases := []struct {
		Input  string
		Result map[string]string
		Err    bool
	}{
		{
			"outputs-empty.tfstate",
			nil,
			false,
		},

		{
			"outputs-basic.tfstate",
			map[string]string{"foo": "bar"},
			false,
		},
	}

	for _, tc := range cases {
		result, err := Outputs(filepath.Join("./test-fixtures", tc.Input))
		if (err != nil) != tc.Err {
			t.Fatalf("bad: %s, %s", tc.Input, err)
		}

		if !reflect.DeepEqual(result, tc.Result) {
			t.Fatalf("bad: %#v, %#v", result, tc.Result)
		}
	}
}
