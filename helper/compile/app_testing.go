package compile

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/hashicorp/otto/otto"
)

var testLock sync.RWMutex
var testOn bool
var testAppOpts *AppOptions

// AppTest enables/disables test mode for the compilation package. When
// test mode is enabled, the test steps can be used to make assertions about
// the compilation process.
//
// Always be sure to defer and disable this.
//
// This should not be used outside of tests. Within tests, this cannot
// be parallelized since it uses global state.
func AppTest(on bool) {
	testLock.Lock()
	defer testLock.Unlock()

	testOn = on
}

// AppTestStepContext is an otto.TestStep that tests the value of something
// in the template context.
type AppTestStepContext struct {
	Key   string
	Value interface{}
}

func (s *AppTestStepContext) Run(c *otto.Core) error {
	testLock.RLock()
	defer testLock.RUnlock()

	if testAppOpts == nil {
		return fmt.Errorf("no context")
	}

	ctx := testAppOpts.Bindata.Context
	if ctx == nil {
		return fmt.Errorf("no context")
	}

	value, ok := ctx[s.Key]
	if !ok || !reflect.DeepEqual(value, s.Value) {
		return fmt.Errorf("bad value for '%s': %#v", s.Key, value)
	}

	return nil
}
