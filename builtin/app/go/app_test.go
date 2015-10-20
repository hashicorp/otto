package goapp

import (
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/otto"
)

func TestApp_impl(t *testing.T) {
	var _ app.App = new(App)
}

func TestApp_dev(t *testing.T) {
	core := otto.TestCoreConfig(t)
	otto.TestAppFixed(t, tuple, config, app)
}
