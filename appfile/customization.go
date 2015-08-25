package appfile

import (
	"strings"
)

// CustomizationSet is a struct that maintains a set of customizations
// from an Appfile and provides helper functions for retrieving and
// filtering them.
//
// Note: While "set" is in the name, this is not a set in the formal
// sense of the word, since customizations can be duplicated.
type CustomizationSet struct {
	// Raw is the raw list of customizations.
	Raw []*Customization
}

// Filter filters the customizations by the given type and returns only
// the matching list of customizations.
func (s *CustomizationSet) Filter(t string) []*Customization {
	// Lowercase the type
	t = strings.ToLower(t)

	// Pre-allocate the result slice to the size of the raw list. There
	// usually aren't that many customization (a handful) so it is easier
	// to just over-allocate here.
	result := make([]*Customization, 0, len(s.Raw))
	for _, c := range s.Raw {
		if c.Type == t {
			result = append(result, c)
		}
	}

	return result
}
