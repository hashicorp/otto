package command

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/otto/infrastructure"
	"github.com/hashicorp/otto/otto"
)

func fixtureDir(n string) string {
	return filepath.Join("./test-fixtures", n)
}

func testCoreConfig(t *testing.T) *otto.CoreConfig {
	return &otto.CoreConfig{
		Infrastructures: map[string]infrastructure.Factory{
			"test": func() (infrastructure.Infrastructure, error) {
				return new(infrastructure.Mock), nil
			},
		},
	}
}
