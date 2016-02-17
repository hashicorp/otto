package hashicorp

import (
	"github.com/hashicorp/otto/foundation"
)

// Tuples is the list of tuples that this built-in foundation implementation knows
// that it can support.
var Tuples = foundation.TupleSlice([]foundation.Tuple{
	{"hashicorp", "aws", "simple"},
})
