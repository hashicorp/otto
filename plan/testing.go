package plan

import (
	"testing"
)

// TestPlan parses a set of plans and fails the test if they do not parse.
func TestPlan(t *testing.T, path string) []*Plan {
	result, err := ParseFile(path)
	if err != nil {
		t.Fatalf("Error parsing plans in %s: %s", path, err)
	}

	return result
}
