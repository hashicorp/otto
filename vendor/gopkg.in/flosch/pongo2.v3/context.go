package pongo2

import (
	"fmt"
	"regexp"
)

var reIdentifiers = regexp.MustCompile("^[a-zA-Z0-9_]+$")

// Use this Context type to provide constants, variables, instances or functions to your template.
//
// pongo2 automatically provides meta-information or functions through the "pongo2"-key.
// Currently, context["pongo2"] contains the following keys:
//  1. version: returns the version string
//
// Template examples for accessing items from your context:
//     {{ myconstant }}
//     {{ myfunc("test", 42) }}
//     {{ user.name }}
//     {{ pongo2.version }}
type Context map[string]interface{}

func (c Context) checkForValidIdentifiers() *Error {
	for k, v := range c {
		if !reIdentifiers.MatchString(k) {
			return &Error{
				Sender:   "checkForValidIdentifiers",
				ErrorMsg: fmt.Sprintf("Context-key '%s' (value: '%+v') is not a valid identifier.", k, v),
			}
		}
	}
	return nil
}

func (c Context) Update(other Context) Context {
	for k, v := range other {
		c[k] = v
	}
	return c
}

// If you're writing a custom tag, your tag's Execute()-function will
// have access to the ExecutionContext. This struct stores anything
// about the current rendering process's Context including
// the Context provided by the user (field Public).
// You can safely use the Private context to provide data to the user's
// template (like a 'forloop'-information). The Shared-context is used
// to share data between tags. All ExecutionContexts share this context.
//
// Please be careful when accessing the Public data.
// PLEASE DO NOT MODIFY THE PUBLIC CONTEXT (read-only).
//
// To create your own execution context within tags, use the
// NewChildExecutionContext(parent) function.
type ExecutionContext struct {
	template *Template

	Autoescape bool
	Public     Context
	Private    Context
	Shared     Context
}

var pongo2MetaContext = Context{
	"version": Version,
}

func newExecutionContext(tpl *Template, ctx Context) *ExecutionContext {
	privateCtx := make(Context)

	// Make the pongo2-related funcs/vars available to the context
	privateCtx["pongo2"] = pongo2MetaContext

	return &ExecutionContext{
		template: tpl,

		Public:     ctx,
		Private:    privateCtx,
		Autoescape: true,
	}
}

func NewChildExecutionContext(parent *ExecutionContext) *ExecutionContext {
	newctx := &ExecutionContext{
		template: parent.template,

		Public:     parent.Public,
		Private:    make(Context),
		Autoescape: parent.Autoescape,
	}
	newctx.Shared = parent.Shared

	// Copy all existing private items
	newctx.Private.Update(parent.Private)

	return newctx
}

func (ctx *ExecutionContext) Error(msg string, token *Token) *Error {
	filename := ctx.template.name
	var line, col int
	if token != nil {
		// No tokens available
		// TODO: Add location (from where?)
		filename = token.Filename
		line = token.Line
		col = token.Col
	}
	return &Error{
		Template: ctx.template,
		Filename: filename,
		Line:     line,
		Column:   col,
		Token:    token,
		Sender:   "execution",
		ErrorMsg: msg,
	}
}

func (ctx *ExecutionContext) Logf(format string, args ...interface{}) {
	ctx.template.set.logf(format, args...)
}
