package plan

import (
	"reflect"
	"testing"
)

func TestTaskArg_Refs(t *testing.T) {
	cases := []struct {
		Value    interface{}
		Expected []string
	}{
		{
			42,
			nil,
		},

		{
			"foo ${bar}",
			[]string{"bar"},
		},

		{
			"foo ${bar+baz}",
			[]string{"bar", "baz"},
		},

		{
			"foo ${bar+bar}",
			[]string{"bar"},
		},
	}

	for _, tc := range cases {
		arg := &TaskArg{Value: tc.Value}
		actual := arg.Refs()
		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("Input: %#v\n\nGot: %#v", tc.Value, actual)
		}
	}
}
