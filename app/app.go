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
	"github.com/hashicorp/otto/ui"
)

// App is the interface that must be implemented by each
// (app type, infra type, infra flavor) 3-tuple.
type App interface {
	// Compile is called to compile the files that are used to manage
	// this application.
	Compile(*Context) (*CompileResult, error)

	// Dev should manage a development environment for this app
	// type. This is called for the local, mutable dev environment
	// where this application is the main thing under development.
	Dev(*Context) error

	// DevDep is called when this app is an upstream dependency
	// of another application that is being developed. This app should
	// build itself for development, and configure the Vagrantfile so
	// that this dependency starts properly on boot.
	//
	// DevDep is given two contexts. The first is the destination
	// app (the one being developed), and the second is the source
	// app (this one that is an upstream dep).
	//
	// The results of this call are cached to speed up development
	// of the destination app until there is a change, which is detected
	// based on VCS.
	DevDep(dst *Context, src *Context) (*DevDep, error)
}

// Context is the context for operations on applications. Some of the
// fields in this struct are only available for certain operations.
type Context struct {
	// Action is the sub-action to take when being executed.
	//
	// ActionArgs is the list of arguments for this action.
	//
	// Both of these fields will only be set for the Execute call.
	Action     string
	ActionArgs []string

	// Dir is the directory that the compilation is allowed to write to
	// for persistant storage of data that is available during task
	// execution. For tasks, this will be the directory that compilation
	// wrote to. Whenever a compilation is done, this directory is
	// cleared. Data that should be persistant across compilations should
	// be stored in the directory service.
	Dir string

	// CacheDir is the directory where data can be cached. This data
	// will persist across compiles of the same version of an Appfile.
	//
	// The App implementation should function under the assumption that
	// this cache directory can be cleared at any time between runs.
	CacheDir string

	// Tuple is the Tuple that identifies this application. This can be
	// used so that an implementatin of App can work with multiple tuple
	// types.
	Tuple Tuple

	// Appfile is the full appfile
	Appfile *appfile.File

	// Application is the application configuration itself from the appfile.
	Application *appfile.Application

	// Ui is the Ui object that can be used to communicate with the user.
	Ui ui.Ui

	// DevDepFragments will be populated with the list of dev dep
	// Vagrantfile fragment paths. This will only be available in the Compile
	// call.
	DevDepFragments []string
}

// CompileResult is the structure containing compilation result values.
type CompileResult struct {
	// DevDepFragmentPath is the path to the Vagrantfile fragment that
	// should be added to other Vagrantfiles when this application is
	// used as a dependency.
	DevDepFragmentPath string
}
