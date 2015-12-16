package scriptpackapp

import (
	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile/detect"
)

// Meta is the metadata for this app type
var Meta = &app.Meta{
	Tuples:    Tuples,
	Detectors: Detectors,
}

// Tuples is the list of tuples that this built-in app implementation knows
// that it can support.
var Tuples = app.TupleSlice([]app.Tuple{
	{"scriptpack", "*", "*"},
})

// Detectors is the list of detectors that trigger this app to be used.
var Detectors = []*detect.Detector{
	&detect.Detector{
		Type: "scriptpack",
		Contents: map[string]string{
			"main.go": `^var ScriptPack =`,
		},
	},
}

// AppFactory is the factory for this app
func AppFactory() app.App {
	return &App{}
}
