package java

import (
	"github.com/hashicorp/otto/builtin/scriptpack/stdlib"
	"github.com/hashicorp/otto/helper/bindata"
	"github.com/hashicorp/otto/scriptpack"
)

//go:generate go generate github.com/hashicorp/otto/builtin/scriptpack/stdlib
//go:generate go-bindata -o=bindata.go -pkg=java -nomemcopy -nometadata ./data/...

// ScriptPack is the exported ScriptPack that can be used.
var ScriptPack = scriptpack.ScriptPack{
	Name: "JAVA",
	Data: bindata.Data{
		Asset:    Asset,
		AssetDir: AssetDir,
	},
	Dependencies: []*scriptpack.ScriptPack{
		&stdlib.ScriptPack,
	},
}
