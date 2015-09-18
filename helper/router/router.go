package router

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"strings"
)

// Router is a helper to route subcommands to specific callbacks.
//
// Actions are available on a lot of commands such as dev, deploy, etc. and
// this can be used to add custom actions.
type Router struct {
	Actions map[string]Action
}

// Action defines an action that is available for the router.
type Action interface {
	// Execute is the callback that'll be called to execute this action.
	Execute(ctx Context) error

	// Help is the help text for this action.
	Help() string

	// Synopsis is the text that will be shown as a short sentence
	// about what this action does.
	Synopsis() string
}

// Context is passed to the router and used to select which action is executed.
// This same value will also be passed down into the selected Action's Execute
// function. This is so that actions typecast the context to access
// implementation-specific data.
type Context interface {
	RouteName() string
	RouteArgs() []string
}

// Route will route the given Context to the proper Action.
func (r *Router) Route(ctx Context) error {
	if _, ok := r.Actions["help"]; !ok {
		r.Actions["help"] = &SimpleAction{
			ExecuteFunc:  r.help,
			SynopsisText: "This help",
		}
	}

	action, ok := r.Actions[ctx.RouteName()]
	if !ok {
		log.Printf("[DEBUG] No action found: %q; executing help.", ctx.RouteName())
		return r.help(ctx)
	}

	return action.Execute(ctx)
}

func (r *Router) help(ctx Context) error {
	// If this is the help command we've been given a specific subcommand
	// to look up, then do that.
	if ctx.RouteName() == "help" && len(ctx.RouteArgs()) > 0 {
		if a, ok := r.Actions[ctx.RouteArgs()[0]]; ok {
			return fmt.Errorf(a.Help())
		}
	}

	// Normal help output...
	var message bytes.Buffer
	if ctx.RouteName() != "" && ctx.RouteName() != "help" {
		message.WriteString(fmt.Sprintf(
			"Unsupported action: %s\n\n", ctx.RouteName()))
	}

	message.WriteString(fmt.Sprintf(
		"The available subcommands are shown below along with a\n" +
			"brief description of what that command does. For more complete\n" +
			"help, call the `help` subcommand with the name of the specific\n" +
			"subcommand you want help for, such as `help foo`.\n\n" +
			"The subcommand '(default)' is the blank subcommand. For this\n" +
			"you don't specify any additional text.\n\n"))

	longestName := len("(default)")
	actionLines := make([]string, 0, len(r.Actions))

	for n, _ := range r.Actions {
		if len(n) > longestName {
			longestName = len(n)
		}
	}
	fmtStr := fmt.Sprintf("    %%%ds\t%%s\n", longestName)

	for n, a := range r.Actions {
		if n == "" {
			n = "(default)"
		}

		actionLines = append(actionLines, fmt.Sprintf(fmtStr, n, a.Synopsis()))
	}

	sort.Strings(actionLines)
	message.WriteString(strings.Join(actionLines, ""))

	return fmt.Errorf(message.String())
}

type SimpleAction struct {
	ExecuteFunc  func(Context) error
	HelpText     string
	SynopsisText string
}

func (sa *SimpleAction) Execute(ctx Context) error {
	return sa.ExecuteFunc(ctx)
}

func (sa *SimpleAction) Help() string {
	return sa.HelpText
}

func (sa *SimpleAction) Synopsis() string {
	return sa.SynopsisText
}
