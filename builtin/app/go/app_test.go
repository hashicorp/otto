package goapp

import (
	"testing"

	"github.com/hashicorp/otto/app"
)

func TestApp_impl(t *testing.T) {
	var _ app.App = new(App)
}
