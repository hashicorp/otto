package otto

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// TestEnvVar must be set to a non-empty value for acceptance tests to run.
const TestEnvVar = "OTTO_ACC"

// TestCase is a single set of tests to run for an app to test the dev
// experience. See Test for more details on how this operates.
type TestCase struct {
	// Precheck, if non-nil, will be called once before the test case
	// runs at all. This can be used for some validation prior to the
	// test running.
	PreCheck func()

	// Core is the core to use for testing. The Test* methods such
	// as TestCoreConfig should be used to create the test core.
	Core *Core

	// Steps are the set of operations that are run for this test case.
	Steps []TestStep

	// Teardown will be called before the test case is over regardless
	// of if the test succeeded or failed. This should return an error
	// in the case that the test can't guarantee all resources were
	// properly cleaned up.
	Teardown TestTeardownFunc
}

// TestStep is a single step within a TestCase.
type TestStep interface {
	// Run is used to run the TestStep. It should return an error if
	// the step failed. If the step fails, then no further steps are
	// called. The Teardown will be called on the TestCase.
	Run(*Core) error
}

// TestTeardownFunc is the callback used for Teardown in TestCase.
type TestTeardownFunc func(*Core) error

// Test performs an acceptance test on a backend with the given test case.
//
// Tests are not run unless an environmental variable "TF_ACC" is
// set to some non-empty value. This is to avoid test cases surprising
// a user by creating real resources.
//
// Tests will fail unless the verbose flag (`go test -v`, or explicitly
// the "-test.v" flag) is set. Because some acceptance tests take quite
// long, we require the verbose flag so users are able to see progress
// output.
func Test(t TestT, c TestCase) {
	// We only run acceptance tests if an env var is set because they're
	// slow and generally require some outside configuration.
	if os.Getenv(TestEnvVar) == "" {
		t.Skip(fmt.Sprintf(
			"Acceptance tests skipped unless env '%s' set",
			TestEnvVar))
		return
	}

	// We require verbose mode so that the user knows what is going on.
	if !testTesting && !testing.Verbose() {
		t.Fatal("Acceptance tests must be run with the -v flag on tests")
		return
	}

	// Run the PreCheck if we have it
	if c.PreCheck != nil {
		c.PreCheck()
	}

	// Check that the core is provided
	if c.Core == nil {
		t.Fatal("Must provide a core")
	}

	// Compile the app
	log.Printf("[WARN] test: compiling appfile...")
	if err := c.Core.Compile(); err != nil {
		t.Fatal("error compiling: ", err)
	}

	// Run the steps
	for i, s := range c.Steps {
		log.Printf("[WARN] Executing test step %d", i+1)
		if err := s.Run(c.Core); err != nil {
			t.Error(fmt.Sprintf("Failed step %d: %s", i+1, err))
			break
		}
	}

	// Cleanup
	if c.Teardown != nil {
		if err := c.Teardown(c.Core); err != nil {
			t.Error(fmt.Sprintf(
				"Teardown failed! Dangling resources may exist. Error:\n\n%s",
				err))
		}
	}
}

// TestT is the interface used to handle the test lifecycle of a test.
//
// Users should just use a *testing.T object, which implements this.
type TestT interface {
	Error(args ...interface{})
	Fatal(args ...interface{})
	Skip(args ...interface{})
}

var testTesting = false
