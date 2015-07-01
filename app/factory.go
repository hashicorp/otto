package app

import (
	"fmt"
	"reflect"
)

// Factory is a factory function for creating infrastructures.
type Factory func() (App, error)

// StructFactory returns a factory function for creating a newly
// instantiated copy of the type of v.
func StructFactory(v App) Factory {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return func() (App, error) {
		raw := reflect.New(t)
		v, ok := raw.Interface().(App)
		if !ok {
			return nil, fmt.Errorf(
				"Failed to instantiate type: %#v", raw.Interface())
		}

		return v, nil
	}
}
