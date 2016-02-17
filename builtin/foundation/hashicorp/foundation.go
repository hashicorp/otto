package hashicorp

import (
	"github.com/hashicorp/otto/foundation"
)

type Foundation struct{}

func (f *Foundation) Compile(*foundation.Context) (*foundation.CompileResult, error) {
	return nil, nil
}

func (f *Foundation) Infra(*foundation.Context) error {
	return nil
}
