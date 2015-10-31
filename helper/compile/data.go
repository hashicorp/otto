package compile

import (
	"github.com/hashicorp/otto/helper/bindata"
)

//go:generate go-bindata -o=bindata.go -pkg=compile -nomemcopy -nometadata ./data/...

// Data is the compiled bindata for this package. This isn't a pointer to
// force a copy so that Context data is never shared.
var Data = bindata.Data{
	Asset:    Asset,
	AssetDir: AssetDir,
}
