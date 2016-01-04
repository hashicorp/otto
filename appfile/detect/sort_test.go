package detect

import (
	"reflect"
	"sort"
	"testing"
)

func TestDetectorList_impl(t *testing.T) {
	var _ sort.Interface = new(DetectorList)
}

func TestDetectorList(t *testing.T) {
	cases := []struct {
		Input  []*Detector
		Output []*Detector
	}{
		{
			Input: []*Detector{
				&Detector{Type: "foo"},
				&Detector{Type: "bar"},
			},
			Output: []*Detector{
				&Detector{Type: "foo"},
				&Detector{Type: "bar"},
			},
		},

		{
			Input: []*Detector{
				&Detector{Type: "foo"},
				&Detector{Type: "bar", Priority: 10},
			},
			Output: []*Detector{
				&Detector{Type: "bar", Priority: 10},
				&Detector{Type: "foo"},
			},
		},
	}

	for i, tc := range cases {
		sort.Sort(DetectorList(tc.Input))
		if !reflect.DeepEqual(tc.Input, tc.Output) {
			t.Fatalf("%d\n\nInput: %#v\n\nOutput: %#v", i, tc.Input, tc.Output)
		}
	}
}
