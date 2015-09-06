package app

import (
	"testing"
)

func TestRouter_default(t *testing.T) {
	var called bool
	executeFunc := func(*Context) error {
		called = true
		return nil
	}

	r := &Router{
		Actions: map[string]*Action{
			"": &Action{
				Execute: executeFunc,
			},
		},
	}

	err := r.Route(&Context{})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !called {
		t.Fatal("should be called")
	}
}

func TestRouter_specific(t *testing.T) {
	var called bool
	executeFunc := func(*Context) error {
		called = true
		return nil
	}

	r := &Router{
		Actions: map[string]*Action{
			"foo": &Action{
				Execute: executeFunc,
			},
		},
	}

	err := r.Route(&Context{Action: "foo"})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !called {
		t.Fatal("should be called")
	}
}
