package app

// Tuple is the tupled used for looking up the App implementation
// for an Appfile. This struct is usually used in its non-pointer form, to be
// a key for maps.
type Tuple struct {
	App         string // App is the app type, i.e. "go"
	Infra       string // Infra is the infra type, i.e. "aws"
	InfraFlavor string // InfraFlavor is the flavor, i.e. "vpc-public-private"
}

// TupleSlice is an alias of []Tuple that implements sort.Interface for
// sorting tuples. See the tests in tuple_test.go to see the sorting order.
type TupleSlice []Tuple

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
