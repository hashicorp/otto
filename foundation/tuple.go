package foundation

import (
	"fmt"
)

// Tuple is the tupled used for looking up the Foundation implementation
// for an Appfile. This struct is usually used in its non-pointer form, to be
// a key for maps.
type Tuple struct {
	Type        string // Type is the foundation type, i.e. "consul"
	Infra       string // Infra is the infra type, i.e. "aws"
	InfraFlavor string // InfraFlavor is the flavor, i.e. "vpc-public-private"
}

func (t Tuple) String() string {
	return fmt.Sprintf("(%q, %q, %q)", t.Type, t.Infra, t.InfraFlavor)
}

// TupleSlice is an alias of []Tuple that implements sort.Interface for
// sorting tuples. See the tests in tuple_test.go to see the sorting order.
type TupleSlice []Tuple

// sort.Interface impl.
func (s TupleSlice) Len() int      { return len(s) }
func (s TupleSlice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s TupleSlice) Less(i, j int) bool {
	if s[i].Type != s[j].Type {
		return s[i].Type < s[j].Type
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

// Lookup looks up a Tuple. This should be used instead of direct [] access
// since it respects wildcards ('*') within the Tuple.
func (m TupleMap) Lookup(t Tuple) Factory {
	// If it just exists, return it
	if f, ok := m[t]; ok {
		return f
	}

	// Eh, this isn't terrible complexity, but we should probably look
	// to do better than this at some point.
	for h, f := range m {
		if h.Type != "*" && h.Type != t.Type {
			continue
		}
		if h.Infra != "*" && h.Infra != t.Infra {
			continue
		}
		if h.InfraFlavor != "*" && h.InfraFlavor != t.InfraFlavor {
			continue
		}

		return f
	}

	return nil
}

// Add is a helper to add another map to this one.
func (m TupleMap) Add(m2 TupleMap) {
	for k, v := range m2 {
		m[k] = v
	}
}
