package otto

import (
	"testing"

	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/ui"
)

// TestCoreConfig returns a CoreConfig that can be used for testing.
func TestCoreConfig(t *testing.T) *CoreConfig {
	return &CoreConfig{
		Infrastructures: map[string]infrastructure.Factory{
			"test": func() (infrastructure.Infrastructure, error) {
				return new(infrastructure.Mock), nil
			},
		},
		Ui: new(ui.Mock),
	}
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
