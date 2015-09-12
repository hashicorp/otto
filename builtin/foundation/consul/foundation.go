package consul

import (
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/helper/compile"
)

//go:generate go-bindata -pkg=consul -nomemcopy -nometadata ./data/...

// Foundation is an implementation of foundation.Foundation
type Foundation struct{}

func (f *Foundation) Compile(ctx *foundation.Context) (*foundation.CompileResult, error) {
	var opts compile.FoundationOptions
	opts = compile.FoundationOptions{
		Ctx: ctx,
		Bindata: &bindata.Data{
			Asset:    Asset,
			AssetDir: AssetDir,
		},
	}

	return compile.Foundation(&opts)
}
