package compile

import (
	"fmt"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/helper/schema"
)

// CustomizationFunc is the callback called for customizations.
type CustomizationFunc func(*schema.FieldData) error

// Customization is used to configure how customizations are processed.
type Customization struct {
	// Schema is the actual schema of the customization configuration. This
	// will be type validated automatically.
	Schema map[string]*schema.FieldSchema

	// Callback is the callback that is called to process this customization.
	// This is guaranteed to be called even if there is no customization set
	// to allow you to setup defaults.
	Callback CustomizationFunc
}

// Merge will merge this customization with the other and return a new
// customization. The original customization is not modified.
func (c *Customization) Merge(other *Customization) *Customization {
	result := &Customization{
		Schema: make(map[string]*schema.FieldSchema),
	}

	// Merge the schemas
	for k, v := range c.Schema {
		result.Schema[k] = v
	}
	for k, v := range other.Schema {
		result.Schema[k] = v
	}

	// Wrap the callbacks
	result.Callback = func(d *schema.FieldData) error {
		if err := c.Callback(d); err != nil {
			return err
		}

		return other.Callback(d)
	}

	return result
}

func processCustomizations(cs *appfile.CustomizationSet, c *Customization) error {
	// We process customizations below by going through multiple
	// passes. We can very likely condense this into one for loop but
	// it helps the semantic understanding to split it out and there should
	// never be so many customizations where the optimizations here matter.

	// We start by going through, building the FieldData.
	rawData := make(map[string]interface{})

	// Go through all the customizations and merge. We only do
	// key-level merging.
	if cs != nil {
		for _, c := range cs.Raw {
			for k, v := range c.Config {
				rawData[k] = v
			}
		}
	}

	// Build the FieldData structure from it
	data := &schema.FieldData{
		Raw:    rawData,
		Schema: c.Schema,
	}

	// Validate it. If it is valid, then we're fine.
	if err := data.Validate(); err != nil {
		return fmt.Errorf("Error in customization: %s", err)
	}

	// Call the callback
	if err := c.Callback(data); err != nil {
		return fmt.Errorf("Error in customization: %s", err)
	}

	return nil
}
