package detect

import (
	"reflect"
	"testing"
)

func TestConfigMerge(t *testing.T) {
	cases := []struct {
		One, Two *Config
		Result   *Config
	}{
		{
			&Config{
				Detectors: []*Detector{
					&Detector{Type: "foo"},
				},
			},
			&Config{
				Detectors: []*Detector{
					&Detector{Type: "bar"},
				},
			},
			&Config{
				Detectors: []*Detector{
					&Detector{Type: "foo"},
					&Detector{Type: "bar"},
				},
			},
		},
	}

	for i, tc := range cases {
		if err := tc.One.Merge(tc.Two); err != nil {
			t.Fatalf("err: %s", err)
		}

		if !reflect.DeepEqual(tc.One, tc.Result) {
			t.Fatalf("bad %d:\n\n%#v", i, tc.One)
		}
	}
}
