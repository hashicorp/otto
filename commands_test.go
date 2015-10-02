package main

import "testing"

func TestDetectors_ordering(t *testing.T) {
	types := []string{}

	for _, detector := range Detectors {
		types = append(types, detector.Type)
	}

	// PHP projects frequently have a package.json
	if !isBefore(types, "php", "node") {
		t.Errorf("php is not before node")
	}

	// Rails projects are also Ruby projects
	if !isBefore(types, "rails", "ruby") {
		t.Errorf("ruby is not before rails")
	}

	// Ruby projects frequently have a package.json
	if !isBefore(types, "ruby", "node") {
		t.Errorf("ruby is not before node")
	}
}

func isBefore(elems []string, a string, b string) bool {
	found_a := false

	for _, elem := range elems {
		if elem == a {
			found_a = true
		}

		if elem == b {
			return found_a
		}
	}

	return false
}
