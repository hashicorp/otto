package otto

import (
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
)

// TestCoreConfig returns a CoreConfig that can be used for testing.
func TestCoreConfig(t *testing.T) *CoreConfig {
	config := &CoreConfig{
		Ui: new(ui.Mock),
	}

	// Add some default mock implementations. These can be overwritten easily
	TestInfra(t, "test", config)
	TestApp(t, app.Tuple{"test", "test", "test"}, config)

	return config
}

// TestInfra adds a mock infrastructure with the given name to the
// core config and returns it.
func TestInfra(t *testing.T, n string, c *CoreConfig) *infrastructure.Mock {
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
func TestApp(t *testing.T, tuple app.Tuple, c *CoreConfig) *app.Mock {
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
func TestFoundation(t *testing.T, tuple foundation.Tuple, c *CoreConfig) *foundation.Mock {
	if c.Foundations == nil {
		c.Foundations = make(map[foundation.Tuple]foundation.Factory)
	}

	result := new(foundation.Mock)
	c.Foundations[tuple] = func() (foundation.Foundation, error) {
		return result, nil
	}

	return result
}
