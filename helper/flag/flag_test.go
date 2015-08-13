package flag

import (
	"flag"
	"reflect"
	"testing"
)

func TestFilterArgs(t *testing.T) {
	cases := []struct {
		Flags []string
		Args  []string
		Inc   []string
		Exc   []string
	}{
		{
			[]string{"foo"},
			[]string{"-foo=bar"},
			[]string{"-foo=bar"},
			[]string{},
		},

		{
			[]string{"foo"},
			[]string{"-foo=bar", "-bar=baz"},
			[]string{"-foo=bar"},
			[]string{"-bar=baz"},
		},

		{
			[]string{"foo"},
			[]string{"hello"},
			[]string{"hello"},
			[]string{"hello"},
		},
	}

	for _, tc := range cases {
		fs := flag.NewFlagSet("", flag.ContinueOnError)
		for _, a := range tc.Flags {
			fs.String(a, "", "")
		}

		inc, exc := FilterArgs(fs, tc.Args)
		if !reflect.DeepEqual(inc, tc.Inc) {
			t.Fatalf(
				"Flags: %#v\n\nArgs: %#v\n\nInc: %#v\n\nActual: %#v",
				tc.Flags, tc.Args, tc.Inc, inc)
		}
		if !reflect.DeepEqual(exc, tc.Exc) {
			t.Fatalf(
				"Flags: %#v\n\nArgs: %#v\n\nExc: %#v\n\nActual: %#v",
				tc.Flags, tc.Args, tc.Exc, exc)
		}
	}
}
