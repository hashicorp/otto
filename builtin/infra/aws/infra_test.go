package aws

import (
	"testing"

	"github.com/hashicorp/otto/infrastructure"
)

func TestInfra_impl(t *testing.T) {
	var _ infrastructure.Infrastructure = new(Infra)
}
