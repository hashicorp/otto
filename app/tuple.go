package app

import (
	"fmt"
)

// Tuple is the tupled used for looking up the App implementation
// for an Appfile. This struct is usually used in its non-pointer form, to be
// a key for maps.
type Tuple struct {
	App         string // App is the app type, i.e. "go"
	Infra       string // Infra is the infra type, i.e. "aws"
	InfraFlavor string // InfraFlavor is the flavor, i.e. "vpc-public-private"
}

func (t Tuple) String() string {
	return fmt.Sprintf("(%q, %q, %q)", t.App, t.Infra, t.InfraFlavor)
}

// TupleSlice is an alias of []Tuple that implements sort.Interface for
// sorting tuples. See the tests in tuple_test.go to see the sorting order.
type TupleSlice []Tuple

// sort.Interface impl.
func (s TupleSlice) Len() int      { return len(s) }
func (s TupleSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TupleSlice) Less(i, j int) bool {
	if s[i].App != s[j].App {
		return s[i].App < s[j].App
	}

	if s[i].Infra != s[j].Infra {
		return s[i].Infra < s[j].Infra
	}

	if s[i].InfraFlavor != s[j].InfraFlavor {
		return s[i].InfraFlavor < s[j].InfraFlavor
	}

	return false
}

// Map turns a TupleSlice into a map where all the tuples in the slice
// are mapped to a single factory function.
func (s TupleSlice) Map(f Factory) TupleMap {
	m := make(TupleMap, len(s))
	for _, t := range s {
		m[t] = f
	}

	return m
}

// TupleMap is an alias of map[Tuple]Factory that adds additional helper
// methods on top to help work with app tuples.
type TupleMap map[Tuple]Factory
