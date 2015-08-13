package aws

import (
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/terraform"
	"github.com/hashicorp/otto/infrastructure"
)

//go:generate go-bindata -pkg=aws -nomemcopy ./data/...

// Infra returns the infrastructure.Infrastructure implementation.
// This function is a infrastructure.Factory.
func Infra() (infrastructure.Infrastructure, error) {
	return &terraform.Infrastructure{
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
		},
	}, nil
}
