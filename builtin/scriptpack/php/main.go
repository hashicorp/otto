package php

import (
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/scriptpack"
)

//go:generate go-bindata -o=bindata.go -pkg=php -nomemcopy -nometadata ./data/...

// ScriptPack is the exported ScriptPack that can be used.
var ScriptPack = scriptpack.ScriptPack{
	Name: "PHP",
	Data: bindata.Data{
		Asset:    Asset,
		AssetDir: AssetDir,
	},
}
