package consul

import(
	"github.com/hashicorp/otto/helper/foundation"
)

// Foundation is an implementation of foundation.Foundation
type Foundation struct{}

func (f *Foundation) Compile(*foundation.Context) (*foundation.CompileResult, error) {
	return nil, nil
}
