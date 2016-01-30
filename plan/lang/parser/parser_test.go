package parser

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/hashicorp/otto/plan"
)

func TestParse(t *testing.T) {
	cases := []struct {
		Name     string
		Err      bool
		Expected []*plan.Plan
	}{
		{
			"empty.ottoplan",
			false,
			nil,
		},

		{
			"empty_single.ottoplan",
			false,
			[]*plan.Plan{
				&plan.Plan{},
			},
		},
	}

	const fixtureDir = "./test-fixtures"
	for _, tc := range cases {
		d, err := ioutil.ReadFile(filepath.Join(fixtureDir, tc.Name))
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		raw, err := Parse(d)
		if (err != nil) != tc.Err {
			t.Fatalf("Input: %s\n\nError: %s", tc.Name, err)
		}

		if p := raw.Plans(); !reflect.DeepEqual(p, tc.Expected) {
			t.Fatalf("Input: %s\n\nBad: %#v", tc.Name, p)
		}
	}
}
