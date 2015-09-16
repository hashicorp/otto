package appfile

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

// Validate validates the Appfile
func (f *File) Validate() error {
	var result error
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

	return result
}
