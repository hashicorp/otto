package rubyapp

import (
	"github.com/hashicorp/otto/app"
)

// Tuples is the list of tuples that this built-in app implementation knows
// that it can support.
var Tuples = app.TupleSlice([]app.Tuple{
	{"rails", "aws", "simple"},
	{"ruby", "aws", "simple"},
	{"ruby", "aws", "vpc-public-private"},
})
