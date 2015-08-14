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
		Pos   []string
	}{
		{
			[]string{"foo"},
			[]string{"-foo=bar"},
			[]string{"-foo=bar"},
			[]string{},
			[]string{},
		},

		{
			[]string{"foo"},
			[]string{"-foo=bar", "-bar=baz"},
			[]string{"-foo=bar"},
			[]string{"-bar=baz"},
			[]string{},
		},

		{
			[]string{"foo"},
			[]string{"-foo", "bar", "-bar=baz"},
			[]string{"-foo", "bar"},
			[]string{"-bar=baz"},
			[]string{},
		},

		{
			[]string{"foo"},
			[]string{"hello"},
			[]string{},
			[]string{},
			[]string{"hello"},
		},

		{
			[]string{"foo"},
			[]string{"-h"},
			[]string{"-h"},
			[]string{},
			[]string{},
		},
	}

	for _, tc := range cases {
		fs := flag.NewFlagSet("", flag.ContinueOnError)
		for _, a := range tc.Flags {
			fs.String(a, "", "")
		}

		inc, exc, pos := FilterArgs(fs, tc.Args)
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
		if !reflect.DeepEqual(pos, tc.Pos) {
			t.Fatalf(
				"Flags: %#v\n\nArgs: %#v\n\nPos: %#v\n\nActual: %#v",
				tc.Flags, tc.Args, tc.Pos, pos)
		}
	}
}
