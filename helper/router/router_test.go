package router

import (
	"flag"
	"testing"

	"github.com/hashicorp/otto/ui"
)

func TestRouter_default(t *testing.T) {
	var called bool
	executeFunc := func(ctx Context) error {
		called = true
		return nil
	}

	r := &Router{
		Actions: map[string]Action{
			"": &SimpleAction{
				ExecuteFunc: executeFunc,
			},
		},
	}

	err := r.Route(&stubContext{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !called {
		t.Fatal("should be called")
	}
}

func TestRouter_specific(t *testing.T) {
	var called bool
	executeFunc := func(ctx Context) error {
		called = true
		return nil
	}

	r := &Router{
		Actions: map[string]Action{
			"foo": &SimpleAction{
				ExecuteFunc: executeFunc,
			},
		},
	}

	err := r.Route(&stubContext{routeName: "foo"})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !called {
		t.Fatal("should be called")
	}
}

func TestRouter_helpErr(t *testing.T) {
	var called bool
	executeFunc := func(ctx Context) error {
		called = true
		return nil
	}

	r := &Router{
		Actions: map[string]Action{
			"help": &SimpleAction{
				ExecuteFunc: executeFunc,
			},

			"foo": &SimpleAction{
				ExecuteFunc: func(Context) error { return ErrHelp },
			},
		},
	}

	err := r.Route(&stubContext{routeName: "foo"})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !called {
		t.Fatal("should be called")
	}
}

func TestRouter_flagHelpErr(t *testing.T) {
	var called bool
	executeFunc := func(ctx Context) error {
		called = true
		return nil
	}

	flagExecuteFunc := func(ctx Context) error {
		fs := flag.NewFlagSet("foo", flag.ContinueOnError)
		return fs.Parse(ctx.RouteArgs())
	}

	r := &Router{
		Actions: map[string]Action{
			"help": &SimpleAction{
				ExecuteFunc: executeFunc,
			},

			"foo": &SimpleAction{
				ExecuteFunc: flagExecuteFunc,
			},
		},
	}

	err := r.Route(&stubContext{
		routeName: "foo",
		routeArgs: []string{"-help"},
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !called {
		t.Fatal("should be called")
	}
}

func TestSimpleAction_impl(t *testing.T) {
	var _ Action = new(SimpleAction)
}

type stubContext struct {
	routeName string
	routeArgs []string
	ui        ui.Ui
}

func (mc *stubContext) RouteName() string {
	return mc.routeName
}

func (mc *stubContext) RouteArgs() []string {
	return mc.routeArgs
}

func (mc *stubContext) UI() ui.Ui {
	return mc.ui
}
