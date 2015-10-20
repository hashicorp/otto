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
