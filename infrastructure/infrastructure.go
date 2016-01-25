package infrastructure

import (
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/plan"
	"github.com/hashicorp/otto/ui"
)

// Infrastructure is an interface that must be implemented by each
// infrastructure type with a method of creating it.
type Infrastructure interface {
	// Creds is called when Otto determines that it needs credentials
	// for this infrastructure provider. The Infra should query the
	// user (or environment) for creds and return them. Otto will
	// handle encrypting, storing, and retrieving the credentials.
	Creds(*Context) (map[string]string, error)

	// VerifyCreds is called with the result of either prompting or
	// retrieving cached credentials. This gives Infrastructure
	// implementations a chance to check that credentials are good before
	// continuing to perform any operations.
	VerifyCreds(*Context) error

	// Compile is called to generate any files that we need.
	Compile(*Context) (*CompileResult, error)

	// Plan is called to plan any changes that are necessary for this
	// infrastructure. This method is expected to potentially make network
	// calls, check for drift, and do whatever is necessary to create the
	// plans. This method should not, however, modify infrastructure
	// in any way.
	Plan(*Context) ([]*plan.Plan, error)
}

// Context is the context for operations on infrastructures. Some of
// the fields in this struct are only available for certain operations.
type Context struct {
	context.Shared

	// Action is the sub-action to take when being executed.
	//
	// ActionArgs is the list of arguments for this action.
	//
	// Both of these fields will only be set for the Execute call.
	Action     string
	ActionArgs []string

	// Dir is the directory that the compilation is allowed to write to
	// for persistant storage of data. For other tasks, this will be the
	// directory that was already populated by compilation.
	Dir string

	// The infrastructure configuration itself from the Appfile. This includes
	// the flavor of the infrastructure we want to launch.
	Infra *appfile.Infrastructure
}

// RouteName implements the router.Context interface so we can use Router
func (c *Context) RouteName() string {
	return c.Action
}

// RouteArgs implements the router.Context interface so we can use Router
func (c *Context) RouteArgs() []string {
	return c.ActionArgs
}

// UI implements router.Context so we can use this in router.Router
func (c *Context) UI() ui.Ui {
	return c.Ui
}

// CompileResult is the structure containing compilation result values.
type CompileResult struct{}
