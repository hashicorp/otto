package vagrant

import (
	"testing"
)

func TestLayeredLayerPaths(t *testing.T) {
	l := &Layered{
		Layers: []*Layer{
			&Layer{
				ID: "foo",
			},
			&Layer{
				ID: "bar",
			},
		},
	}

	// Not a great way to test this, but at the very least we can check that
	// the number is valid and that we have paths for everything.
	paths := l.LayerPaths()
	if len(paths) != len(l.Layers) {
		t.Fatalf("bad: %#v", paths)
	}
	for _, l := range l.Layers {
		if _, ok := paths[l.ID]; !ok {
			t.Fatalf("bad: %#v", paths)
		}
	}
}
