package app

import (
	"reflect"
	"sort"
	"testing"
)

func TestTupleSlice_sort(t *testing.T) {
	tuples := []Tuple{
		{"go", "aws", "b"},
		{"rails", "aws", "b"},
		{"go", "google", "a"},
		{"go", "aws", "a"},
	}

	expected := []Tuple{
		{"go", "aws", "a"},
		{"go", "aws", "b"},
		{"go", "google", "a"},
		{"rails", "aws", "b"},
	}

	sort.Sort(TupleSlice(tuples))
	if !reflect.DeepEqual(tuples, expected) {
		t.Fatalf("bad: %#v", tuples)
	}
}

func TestTupleMap_Lookup(t *testing.T) {
	var value int
	factory := func(v int) Factory {
		return func() (App, error) {
			value = v
			return nil, nil
		}
	}

	cases := []struct {
		M        TupleMap
		T        Tuple
		Expected int
	}{
		{
			map[Tuple]Factory{
				Tuple{"foo", "bar", "baz"}: factory(7),
			},
			Tuple{"foo", "bar", "baz"},
			7,
		},

		{
			map[Tuple]Factory{
				Tuple{"foo", "*", "baz"}: factory(7),
			},
			Tuple{"foo", "bar", "baz"},
			7,
		},

		{
			map[Tuple]Factory{
				Tuple{"foo", "bar", "*"}: factory(7),
			},
			Tuple{"foo", "bar", "baz"},
			7,
		},

		{
			map[Tuple]Factory{
				Tuple{"foo", "*", "*"}: factory(7),
			},
			Tuple{"foo", "bar", "baz"},
			7,
		},

		{
			map[Tuple]Factory{
				Tuple{"foo", "*", "*"}:     factory(12),
				Tuple{"foo", "bar", "baz"}: factory(7),
			},
			Tuple{"foo", "bar", "baz"},
			7,
		},
	}

	for i, tc := range cases {
		// Reset the value
		value = 0

		f := tc.M.Lookup(tc.T)
		if f != nil {
			f()
		}

		if value != tc.Expected {
			t.Fatalf("%d: bad: %d", i, value)
		}
	}
}
