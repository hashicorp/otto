package nodeapp

import (
	"path/filepath"
	"testing"

	"github.com/hashicorp/otto/helper/vagrant"
	"github.com/hashicorp/otto/otto"
)

func TestApp_dev(t *testing.T) {
	otto.Test(t, otto.TestCase{
		Core: otto.TestCore(t, &otto.TestCoreOpts{
			Path: filepath.Join("./test-fixtures", "basic", "Appfile"),
			App:  new(App),
		}),

		Steps: []otto.TestStep{
			&vagrant.DevTestStepInit{},

			// Verify we have Node
			&vagrant.DevTestStepGuestScript{
				Command: "node --version",
			},
		},

		Teardown: vagrant.DevTestTeardown,
	})
}
