package rubyapp

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

			// Verify we have Ruby
			&vagrant.DevTestStepGuestScript{
				Command: "ruby --version | grep '2.2'",
			},
			&vagrant.DevTestStepGuestScript{
				Command: "bundle --version",
			},

			// Verify everything works
			&vagrant.DevTestStepGuestScript{
				Command: "bundle exec ruby app.rb | grep hello",
			},
		},

		Teardown: vagrant.DevTestTeardown,
	})
}
