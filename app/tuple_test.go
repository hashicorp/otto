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
