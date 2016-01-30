package parser

import (
	"fmt"

	"github.com/hashicorp/otto/plan/lang/token"
)

// PosError is a parse error that contains a position.
type PosError struct {
	Pos token.Pos
	Err error
}

func (e *PosError) Error() string {
	return fmt.Sprintf("At %s: %s", e.Pos, e.Err)
}
