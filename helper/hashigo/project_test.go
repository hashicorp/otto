package hashigo

import (
	"testing"
)

func TestVersionRe(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		{
			"0.8.7",
			"0.8.7",
		},

		{
			"Terraform v0.6.4-dev (4d37704d532ae3effdfe2c5b6254bee0b94e8d8e+CHANGES)\n",
			"0.6.4-dev",
		},

		{
			"Consul v0.5.0\nConsul Protocol: 2 (Understands back to: 1)",
			"0.5.0",
		},
	}

	for _, tc := range cases {
		matches := versionRe.FindStringSubmatch(tc.Input)
		if len(matches) == 0 {
			t.Fatalf("bad: %s", tc.Input)
		}
		if matches[1] != tc.Output {
			t.Fatalf("bad: %s != %s\n\n%s", matches[1], tc.Output, tc.Input)
		}
	}
}
