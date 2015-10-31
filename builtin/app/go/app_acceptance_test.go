package goapp

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

			// Verify we have Go
			&vagrant.DevTestStepGuestScript{
				Command: "go version",
			},

			// Verify we can build immediately (we should be in the directory)
			&vagrant.DevTestStepGuestScript{
				Command: "grep '42' <<< $(go build -o test-output && ./test-output 2>&1)",
			},
		},

		Teardown: vagrant.DevTestTeardown,
	})
}
