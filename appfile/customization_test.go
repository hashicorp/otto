package appfile

import (
	"reflect"
	"testing"
)

func TestCustomizationSetFilter(t *testing.T) {
	cases := []struct {
		Raw        []*Customization
		FilterType string
		Result     []*Customization
	}{
		{
			[]*Customization{
				&Customization{Type: "foo"},
				&Customization{Type: "bar"},
			},
			"foo",
			[]*Customization{
				&Customization{Type: "foo"},
			},
		},

		{
			[]*Customization{
				&Customization{Type: "foo"},
				&Customization{Type: "bar"},
			},
			"fOo",
			[]*Customization{
				&Customization{Type: "foo"},
			},
		},
	}

	for i, tc := range cases {
		set := &CustomizationSet{Raw: tc.Raw}
		actual := set.Filter(tc.FilterType)
		if !reflect.DeepEqual(actual, tc.Result) {
			t.Fatalf("bad %d:\n\n%#v\n\n%#v", i, actual, tc.Result)
		}
	}
}
