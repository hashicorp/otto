package goapp

import (
	"github.com/hashicorp/otto/app"
)

// App is an implementation of app.App
type App struct{}

func (a *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	return nil, nil
}
