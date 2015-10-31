package otto

import (
	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/infrastructure"
)

// TestInfra adds a mock infrastructure with the given name to the
// core config and returns it.
func TestInfra(t TestT, n string, c *CoreConfig) *infrastructure.Mock {
	if c.Infrastructures == nil {
		c.Infrastructures = make(map[string]infrastructure.Factory)
	}

	result := new(infrastructure.Mock)
	c.Infrastructures[n] = func() (infrastructure.Infrastructure, error) {
		return result, nil
	}

	return result
}

// TestApp adds a mock app with the given tuple to the core config.
func TestApp(t TestT, tuple app.Tuple, c *CoreConfig) *app.Mock {
	if c.Apps == nil {
		c.Apps = make(map[app.Tuple]app.Factory)
	}

	result := new(app.Mock)
	c.Apps[tuple] = func() (app.App, error) {
		return result, nil
	}

	return result
}

// TestFoundation adds a mock foundation with the given tuple to the core config.
func TestFoundation(t TestT, tuple foundation.Tuple, c *CoreConfig) *foundation.Mock {
	if c.Foundations == nil {
		c.Foundations = make(map[foundation.Tuple]foundation.Factory)
	}

	result := new(foundation.Mock)
	c.Foundations[tuple] = func() (foundation.Foundation, error) {
		return result, nil
	}

	return result
}
