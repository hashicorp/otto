package compile

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/schema"
)

// Customization defines how customizations are handled during
// compilation and are used for both App and Infra compilations.
//
// Customizations are the "customization" stanzas within the Appfile.
//
// Customizations are processed by querying the type given and then
// calling the Callback. The CustomizationResult that is returned will
// modify the behavior of the compilation process.
type Customization struct {
	// Type is the type of the customization, such as "dev"
	Type string

	// Schema is the schema for the data. This will be automatically
	// validated with the data from the configuration.
	Schema map[string]*schema.FieldSchema

	// Callback is called to process this customization.
	Callback CustomizationFunc
}

// CustomizationFunc is the callback called for customizations.
type CustomizationFunc func(*schema.FieldData) error

type processOpts struct {
	Customizations []*Customization

	Appfile *appfile.File
	Bindata *bindata.Data
}

func processCustomizations(opts *processOpts) error {
	// We process customizations below by going through multiple
	// passes. We can very likely condense this into one for loop but
	// it helps the semantic understanding to split it out and there should
	// never be so many customizations where the optimizations here matter.

	// We start by going through, building the FieldData.
	data := make([]*schema.FieldData, len(opts.Customizations))
	for i, c := range opts.Customizations {
		raw := make(map[string]interface{})

		// Grab the real customizations
		cs := opts.Appfile.Customization.Filter(c.Type)
		if len(cs) > 0 {
			// We just want the last one. We don't do any merging for now
			// or validation of the earlier ones. I'm sure this will cause problems
			// one day.
			realC := cs[len(cs)-1]
			raw = realC.Config
		}

		// Build the FieldData structure from it
		data[i] = &schema.FieldData{
			Raw:    raw,
			Schema: c.Schema,
		}
	}

	// Validate all the field data
	var err error
	for i, d := range data {
		// This is a sparse slice, so if its nil ignore it
		if d == nil {
			continue
		}

		// Validate it. If it is valid, then we're fine.
		verr := d.Validate()
		if verr == nil {
			continue
		}

		// Invalid, record the error
		c := opts.Customizations[i]
		err = multierror.Append(err, fmt.Errorf(
			"Error in '%s' customization: %s", c.Type, verr))
	}

	// If we have validation errors, return now
	if err != nil {
		return err
	}

	// Go through the fields, call the callbacks, and record those results
	for i, d := range data {
		if d == nil {
			continue
		}

		c := opts.Customizations[i]
		if cerr := c.Callback(d); cerr != nil {
			err = multierror.Append(err, fmt.Errorf(
				"Error in '%s' customization: %s", c.Type, cerr))
			continue
		}
	}

	return err
}
