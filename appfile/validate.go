package appfile

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

// Validate validates the Appfile
func (f *File) Validate() error {
	var result error

	// Basic checking for stanzas
	if f.Application == nil {
		result = multierror.Append(result, fmt.Errorf(
			"'application' stanza required for Appfile"))
	}
	if f.Project == nil {
		result = multierror.Append(result, fmt.Errorf(
			"'project' stanza required for Appfile"))
	}
	if f.Infrastructure == nil {
		result = multierror.Append(result, fmt.Errorf(
			"'infra' stanza required for Appfile"))
	}

	// Verify the application itself
	if f.Application != nil {
		if f.Application.Name == "" {
			result = multierror.Append(result, fmt.Errorf(
				"application: name is required"))
		}
		if f.Application.Type == "" {
			result = multierror.Append(result, fmt.Errorf(
				"application: type is required"))
		}
	}

	return result
}
