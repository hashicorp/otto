package app

import (
	"bytes"
	"fmt"
)

// Router is a helper to route actions on the Context to specific callbacks.
//
// Actions are available on a lot of commands such as dev, deploy, etc. and
// this can be used to add custom actions.
type Router struct {
	Actions map[string]*Action
}

// Action defines an action that is available for the router.
type Action struct {
	// Execute is the callback that'll be called to execute this action.
	Execute ActionFunc

	// Help is the help text for this action.
	Help string

	// Synopsis is the text that will be shown as a short sentence
	// about what this action does.
	Synopsis string
}

// ActionFunc is the callback for the router.
type ActionFunc func(*Context) error

// Route will route the given Context to the proper Action.
func (r *Router) Route(ctx *Context) error {
	if _, ok := r.Actions["help"]; !ok {
		r.Actions["help"] = &Action{
			Execute:  r.help,
			Synopsis: "This help",
		}
	}

	action, ok := r.Actions[ctx.Action]
	if !ok {
		return r.help(ctx)
	}

	return action.Execute(ctx)
}

func (r *Router) help(ctx *Context) error {
	// If this is the help command we've been given a specific subcommand
	// to look up, then do that.
	if ctx.Action == "help" && len(ctx.ActionArgs) > 0 {
		if a, ok := r.Actions[ctx.ActionArgs[0]]; ok {
			return fmt.Errorf(a.Help)
		}
	}

	// Normal help output...
	var message bytes.Buffer
	if ctx.Action != "" && ctx.Action != "help" {
		message.WriteString(fmt.Sprintf(
			"Unsupported action: %s\n\n", ctx.Action))
	}

	message.WriteString(fmt.Sprintf(
		"The available subcommands are shown below along with a\n" +
			"brief description of what that command does. For more complete\n" +
			"help, call the `help` subcommand with the name of the specific\n" +
			"subcommand you want help for, such as `help foo`.\n\n" +
			"The subcommand '<default>' is the blank subcommand. For this\n" +
			"you don't specify any additional text.\n\n"))

	// TODO: sort, spacing
	for n, a := range r.Actions {
		if n == "" {
			n = "<default>"
		}

		message.WriteString(fmt.Sprintf(
			"    %s\t\t%s\n", n, a.Synopsis))
	}

	return fmt.Errorf(message.String())
}
