package directory

import (
	"github.com/hashicorp/otto/helper/uuid"
)

// Dev represents a development environment of an app.
type Dev struct {
	// Lookup information for the Deploy. AppID is all that is required.
	Lookup

	// These fields should be set for Put and will be populated on Get
	State DevState // State of the dev environment

	// Private fields. These are usually set on Get or Put.
	//
	// DO NOT MODIFY THESE.
	ID string
}

func (d *Dev) IsReady() bool {
	return d != nil && d.State == DevStateReady
}

func (d *Dev) MarkReady() {
	d.State = DevStateReady
}

func (d *Dev) setId() {
	d.ID = uuid.GenerateUUID()
}
