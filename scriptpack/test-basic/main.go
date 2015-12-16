package testBasic

import (
	"github.com/hashicorp/otto/helper/bindata"
)

//go:generate go-bindata -o=bindata.go -pkg=testBasic -nomemcopy -nometadata ./data/...

var Bindata = bindata.Data{
	Asset:    Asset,
	AssetDir: AssetDir,
}
