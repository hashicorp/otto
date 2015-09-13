package foundation

import (
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/context"
)

// Foundation is the interface that must be implemented by each
// foundation. A foundation is a fundamental building block of a
// real infrastructure, and can be categorized such as service discovery,
// security, etc.
//
// Foundations are bound to a (name, infra type, infra flavor) 3-tuple.
type Foundation interface {
	// Compile is called to compile the files that are used to manage
	// this foundation.
	Compile(*Context) (*CompileResult, error)

	// Infra is called to build or destroy the infrastructure for this
	// foundation. The "Action" field in the Context can be used to
	// determine the desired action. This will be either "" (build)
	// or "destroy". Foundations currently don't support any other
	// actions.
	Infra(*Context) error
}

// Context is the context for operations on a Foundation.
type Context struct {
	context.Shared

	// Action is the sub-action to take when being executed.
	//
	// ActionArgs is the list of arguments for this action.
	//
	// Both of these fields will only be set for the Infra call currently.
	Action     string
	ActionArgs []string

	// Config is the raw configuration from the Appfile itself for
	// this foundation.
	Config map[string]interface{}

	// AppConfig is the foundation configuration that was returned by the
	// application that we're working with. This is only available during
	// the Compile function if we're compiling for an application.
	//
	// It should be expected during compilation that this might be nil.
	// The cases where it is nil are not currently well defined, but the
	// behavior in the nil case should be to do nothing except Deploy.
	AppConfig *Config

	// Dir is the directory that the compilation is allowed to write to
	// for persistant storage of data that is available during task
	// execution. For tasks, this will be the directory that compilation
	// wrote to. Whenever a compilation is done, this directory is
	// cleared. Data that should be persistant across compilations should
	// be stored in the directory service.
	Dir string

	// Appfile is the full appfile
	Appfile *appfile.File

	// Tuple is the tuple used for this foundation.
	Tuple Tuple
}

// CompileResult is the structure containing compilation result values.
//
// This is empty now but may be used in the future.
type CompileResult struct{}
