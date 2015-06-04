package infrastructure

import (
	"fmt"
	"reflect"
)

// Factory is a factory function for creating infrastructures.
type Factory func() (Infrastructure, error)

// StructFactory returns a factory function for creating a newly
// instantiated copy of the type of v.
func StructFactory(v Infrastructure) Factory {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return func() (Infrastructure, error) {
		raw := reflect.New(t)
		v, ok := raw.Interface().(Infrastructure)
		if !ok {
			return nil, fmt.Errorf(
				"Failed to instantiate type: %#v", raw.Interface())
		}

		return v, nil
	}
}
