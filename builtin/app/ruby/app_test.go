package rubyapp

import (
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
)

func TestApp_impl(t *testing.T) {
	var _ app.App = new(App)
}

func TestAppImplicit(t *testing.T) {
	cases := []struct {
		Dir  string
		Deps []string
	}{
		{
			"implicit-none",
			nil,
		},

		{
			"implicit-redis",
			[]string{"github.com/hashicorp/otto/examples/redis"},
		},
	}

	for _, tc := range cases {
		errPrefix := fmt.Sprintf("In '%s': ", tc.Dir)

		// Build our context we send in
		var ctx app.Context
		ctx.Appfile = &appfile.File{Path: filepath.Join("./test-fixtures", tc.Dir, "Appfile")}

		// Get the implicit file
		var a App
		f, err := a.Implicit(&ctx)
		if err != nil {
			t.Fatalf("%s: %s", errPrefix, err)
		}
		if (len(tc.Deps) == 0) != (f == nil) {
			// Complicated statement above but basically: should be nil if
			// we expected no deps, and should not be nil if we expect deps
			t.Fatalf("%s: deps: %#v\n\ninvalid file: %#v", errPrefix, tc.Deps, f)
		}
		if f == nil {
			continue
		}

		// Build the deps we got and sort them for determinism
		actual := make([]string, 0, len(f.Application.Dependencies))
		for _, dep := range f.Application.Dependencies {
			actual = append(actual, dep.Source)
		}
		sort.Strings(actual)
		sort.Strings(tc.Deps)

		// Test
		if !reflect.DeepEqual(actual, tc.Deps) {
			t.Fatalf("%s\n\ngot: %#v\n\nexpected: %#v", errPrefix, actual, tc.Deps)
		}
	}
}
