package plan

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	cases := []struct {
		Name     string
		Err      bool
		Expected []*Plan
	}{
		{
			"empty.hcl",
			false,
			nil,
		},

		{
			"basic.hcl",
			false,
			[]*Plan{
				&Plan{
					Description: "foo",
				},
			},
		},

		{
			"basic-tasks.hcl",
			false,
			[]*Plan{
				&Plan{
					Tasks: []*Task{
						&Task{
							Type:        "foo",
							Description: "desc foo",
							Args: map[string]*TaskArg{
								"foo": &TaskArg{Value: "bar"},
							},
						},

						&Task{
							Type:        "bar",
							Description: "desc bar",
						},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		path := filepath.Join("./test-fixtures", tc.Name)
		actual, err := ParseFile(path)
		if (err != nil) != tc.Err {
			t.Fatalf("%s: err: %s", tc.Name, err)
		}
		if err != nil {
			continue
		}

		if !reflect.DeepEqual(actual, tc.Expected) {
			t.Fatalf("%s: bad: %#v", tc.Name, actual)
		}
	}
}
