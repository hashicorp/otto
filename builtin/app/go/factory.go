package goapp

import (
	"github.com/hashicorp/otto/app"
)

func AppFactory() app.App {
	return &App{}
}
