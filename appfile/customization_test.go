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
				&Customization{Name: "foo"},
				&Customization{Name: "bar"},
			},
			"foo",
			[]*Customization{
				&Customization{Name: "foo"},
			},
		},

		{
			[]*Customization{
				&Customization{Name: "foo"},
				&Customization{Name: "bar"},
			},
			"fOo",
			[]*Customization{
				&Customization{Name: "foo"},
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
