package foundation

import (
	"fmt"
	"reflect"
)

// Factory is a factory function for creating foundations.
type Factory func() (Foundation, error)

// StructFactory returns a factory function for creating a newly
// instantiated copy of the type of v.
func StructFactory(v Foundation) Factory {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return func() (Foundation, error) {
		raw := reflect.New(t)
		v, ok := raw.Interface().(Foundation)
		if !ok {
			return nil, fmt.Errorf(
				"Failed to instantiate type: %#v", raw.Interface())
		}

		return v, nil
	}
}
