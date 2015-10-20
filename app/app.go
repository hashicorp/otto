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
	"github.com/hashicorp/otto/context"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/ui"
)

// App is the interface that must be implemented by each
// (app type, infra type, infra flavor) 3-tuple.
type App interface {
	// Compile is called to compile the files that are used to manage
	// this application.
	Compile(*Context) (*CompileResult, error)

	// Build is called to build the deployable artifact for this
	// application.
	Build(*Context) error

	// Deploy is called to deploy this application. The existence of
	// a prior build artifact is confirmed before this is called.
	Deploy(*Context) error

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
	//
	// The resulting DevDep can be nil if nothing needs to be done that
	// is part of the DevDep structure. Any DevDepFragments from the
	// compilation will still be used, of course.
	DevDep(dst *Context, src *Context) (*DevDep, error)
}

// Context is the context for operations on applications. Some of the
// fields in this struct are only available for certain operations.
type Context struct {
	context.Shared

	// CompileResult is the result of the compilation. This is set on
	// all calls except Compile to be the data from the compilation. This
	// can be used to check compile versions, for example.
	CompileResult *CompileResult

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

	// LocalDir is the directory where data local to this single Appfile
	// will be stored; it isn't cleared for compilation.
	LocalDir string

	// Tuple is the Tuple that identifies this application. This can be
	// used so that an implementatin of App can work with multiple tuple
	// types.
	Tuple Tuple

	// Application is the application configuration itself from the appfile.
	Application *appfile.Application

	// DevDepFragments will be populated with the list of dev dep
	// Vagrantfile fragment paths. This will only be available in the Compile
	// call.
	DevDepFragments []string

	// DevIPAddress is a local IP address in the private address space
	// that can be used for a development environment. Otto core
	// does its best to ensure this is unused.
	//
	// This is only available if this app is the root application being
	// developed (dependencies don't get an IP).
	DevIPAddress string
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
type CompileResult struct {
	// Version is the version of the compiled result. This is purely metadata:
	// the app itself should use this to detect certain behaviors on run.
	Version uint32

	// FoundationConfig is the configuration for the various foundational
	// elements of Otto.
	FoundationConfig foundation.Config

	// DevDepFragmentPath is the path to the Vagrantfile fragment that
	// should be added to other Vagrantfiles when this application is
	// used as a dependency.
	DevDepFragmentPath string

	// FoundationResults are the compilation results of the foundations.
	//
	// This is populated by Otto core and any set value here will be ignored.
	FoundationResults []*foundation.CompileResult
}
