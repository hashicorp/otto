package directory

import (
	"reflect"
	"sort"
	"testing"
)

func TestInfraSlice_impl(t *testing.T) {
	var _ sort.Interface = new(InfraSlice)
}

func TestInfraSlice_sort(t *testing.T) {
	cases := []struct {
		Desc     string
		Input    []*Infra
		Expected []*Infra
	}{
		{
			"basic",
			[]*Infra{
				&Infra{Name: "foo"},
				&Infra{Name: "bar"},
				&Infra{Name: "quux"},
			},
			[]*Infra{
				&Infra{Name: "bar"},
				&Infra{Name: "foo"},
				&Infra{Name: "quux"},
			},
		},
	}

	for _, tc := range cases {
		sort.Sort(InfraSlice(tc.Input))
		if !reflect.DeepEqual(tc.Expected, tc.Input) {
			t.Fatalf("%s bad: %#v", tc.Desc, tc.Input)
		}
	}
}
