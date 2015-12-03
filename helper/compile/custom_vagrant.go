package compile

import (
	"github.com/hashicorp/otto/helper/schema"
)

// VagrantCustomizations returns common Vagrant customizations that work
// with the default settings of the app compilation helper.
func VagrantCustomizations(opts *AppOptions) *Customization {
	return &Customization{
		Callback: vagrantCustomizationCallback(opts),
		Schema: map[string]*schema.FieldSchema{
			"vagrantfile": &schema.FieldSchema{
				Type:        schema.TypeString,
				Description: "Vagrantfile contents to append for development.",
			},
		},
	}
}

func vagrantCustomizationCallback(opts *AppOptions) CustomizationFunc {
	return func(d *schema.FieldData) error {
		opts.Bindata.Context["dev_extra_vagrantfile"] = d.Get("vagrantfile")
		return nil
	}
}
