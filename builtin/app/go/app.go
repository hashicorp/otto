package goapp

import (
	"fmt"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/helper/bindata"
)

//go:generate go-bindata -pkg=goapp -nomemcopy ./data/...

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	data := &bindata.Data{
		Asset:    Asset,
		AssetDir: AssetDir,
	}

	prefix := fmt.Sprintf("data/%s-%s", ctx.Tuple.Infra, ctx.Tuple.InfraFlavor)
	if err := data.CopyDir(ctx.Dir, prefix); err != nil {
		return nil, err
	}

	return nil, nil
}
