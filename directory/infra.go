package directory

import (
	"github.com/hashicorp/otto/helper/uuid"
)

// Infra represents the data stored in the directory service about
// Infrastructures.
type Infra struct {
	// Lookup information for the Infra. The only required field for
	// this is Infra. Optionally you may also specify Foundation to
	// get the infrastructure data for a foundation.
	Lookup

	// State is the state of this infrastructure. This is important since
	// it is possible for there to be a partial state. If we're in a
	// partial state then deploys and such can't go through yet.
	State InfraState

	// Outputs are the output data from the infrastructure step.
	// This is an opaque blob that is dependent on each infrastructure
	// type. Please refer to docs of a specific infra to learn more about
	// what values are here.
	Outputs map[string]string `json:"outputs"`

	// Private fields. These are usually set on Get or Put.
	//
	// DO NOT MODIFY THESE.
	ID string
}

func (i *Infra) IsPartial() bool {
	return i != nil && i.State == InfraStatePartial
}

func (i *Infra) IsReady() bool {
	return i != nil && i.State == InfraStateReady
}

func (i *Infra) setId() {
	i.ID = uuid.GenerateUUID()
}
