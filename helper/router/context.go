package router

import (
	"github.com/hashicorp/otto/ui"
)

// simpleContext is an implementation of Context for returning fixed values.
type simpleContext struct {
	Name  string
	Args  []string
	UIVal ui.Ui
}

func (c *simpleContext) RouteName() string   { return c.Name }
func (c *simpleContext) RouteArgs() []string { return c.Args }
func (c *simpleContext) UI() ui.Ui           { return c.UIVal }
