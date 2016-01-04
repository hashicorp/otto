package nodeapp

import (
	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile/detect"
)

// AppFactory is the factory for this app
func AppFactory() app.App {
	return &App{}
}

// Meta is the metadata for this app type
var Meta = &app.Meta{
	Tuples:    Tuples,
	Detectors: Detectors,
}

// Tuples is the list of tuples that this built-in app implementation knows
// that it can support.
var Tuples = app.TupleSlice([]app.Tuple{
	{"node", "aws", "simple"},
})

// Detectors is the list of detectors that trigger this app to be used.
var Detectors = []*detect.Detector{
	&detect.Detector{
		Type: "node",
		File: []string{"package.json"},

		// Slightly lower priority since many web frameworks can contain
		// a package.json these days while being written in another
		// language.
		Priority: -1,
	},
}
