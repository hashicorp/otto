package packer

import (
	"reflect"
	"testing"
)

func TestParseArtifactAws(t *testing.T) {
	cases := map[string]*Output {
		"aws":
			&Output{
				Timestamp: "1376289459",
				Target:    "",
				Type:      "ui",
				Data: []string{
					"say",
					"id",
					"region:aws_id",
				},
			},
	}

	for name, output := range cases {
		result := make(map[string]string)
		ParseArtifact(result)(output)
		if !reflect.DeepEqual(result["region"], "aws_id") {
			t.Fatalf("%s\n\n%#v\n\n%#v", name, "aws_id", result["region"])
		}
	}
}

func TestParseArtifactDigitalocean(t *testing.T) {
	cases := map[string]*Output {
		"digitalocean":
			&Output{
				Timestamp: "1376289459",
				Target:    "",
				Type:      "ui",
				Data: []string{
					"say",
					"id",
					"digitalocean_id",
				},
			},
	}

	for name, output := range cases {
		result := make(map[string]string)
		ParseArtifact(result)(output)
		if !reflect.DeepEqual(result["sfo1"], "digitalocean_id") {
			t.Fatalf("%s\n\n%#v\n\n%#v", name, "digitalocean_id", result["sfo1"])
		}
	}
}
