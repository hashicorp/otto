// app contains the interfaces and structures for application type
// implementations for Otto. Applications are the components that
// know how to dev, build, and deploy a certain kind of application such
// as Rails, PHP, etc.
//
// All app implementations are built specific to a certain 3-tuple:
// (app type, infra type, infra flavor). For example:
// (rails, aws, vpc-public-private). The app implementation then only
// needs to know how to satisfy that specific 3-tuple.
//
// When building app plugins, it is possible for that plugin to support
// multiple matrix elements, but each implementation of the interface
// is expeced to only implement one.
package app

import (
	"github.com/hashicorp/otto/appfile"
)

// App is the interface that must be implemented by each
// (app type, infra type, infra flavor) 3-tuple.
type App interface {
	Compile(*Context) (*CompileResult, error)
}

// Context is the context for operations on applications. Some of the
// fields in this struct are only available for certain operations.
type Context struct {
	// Dir is the directory that the compilation is allowed to write to
	// for persistant storage of data that is available during task
	// execution. For tasks, this will be the directory that compilation
	// wrote to. Whenever a compilation is done, this directory is
	// cleared. Data that should be persistant across compilations should
	// be stored in the directory service.
	Dir string

	// Tuple is the Tuple that identifies this application. This can be
	// used so that an implementatin of App can work with multiple tuple
	// types.
	Tuple Tuple

	// Application is the application configuration itself from the appfile.
	Application *appfile.Application
}

// CompileResult is the structure containing compilation result values.
type CompileResult struct{}
