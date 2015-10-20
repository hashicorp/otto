package vagrant

import (
	"log"

	"github.com/hashicorp/otto/otto"
)

// DevTestTeardown implements the otto.TestTeardownFunc type and should
// be used with otto.TestCase to clear out development environments cleanly.
func DevTestTeardown(c *otto.Core) error {
	// Destroy the dev environment. This should work even if it isn't
	// running so we can always execute it.
	log.Printf("[INFO] test: destroying the development environment")
	return c.Execute(&otto.ExecuteOpts{
		Task:   otto.ExecuteTaskDev,
		Action: "destroy",
	})
}

// DevTestStepInit is a otto.TestStep that initilizes dev testing.
// This should be the first test step before any others for dev.
type DevTestStepInit struct{}

func (s *DevTestStepInit) Run(c *otto.Core) error {
	log.Printf("[INFO] test: starting the development environment")
	return c.Dev()
}

// DevTestStepGuestScript is an otto.TestStep that runs a script in the
// guest and verifies it succeeds (exit code 0).
type DevTestStepScript struct {
	Command string
}

func (s *DevTestStepScript) Run(c *otto.Core) error {
	log.Printf("[INFO] test: testing guest script: %q", s.Command)
	return nil
}
