package directory

import (
	"github.com/hashicorp/otto/helper/uuid"
)

// DeployState is used to track the state of a deploy.
//
// This is required because a deploy is entered in the directory
// prior to the deploy actually happening so that we can always look
// up any binary blobs stored with a deploy even if it fails.
type DeployState byte

const (
	DeployStateInvalid DeployState = 0
	DeployStateNew     DeployState = iota
	DeployStateFail
	DeployStateSuccess
)

// Deploy represents a deploy of an App.
type Deploy struct {
	App         string // App is the app type, i.e. "go"
	Infra       string // Infra is the infra type, i.e. "aws"
	InfraFlavor string // InfraFlavor is the flavor, i.e. "vpc-public-private"

	// These fields should be set for Put and will be populated on Get
	State  DeployState       // State of the deploy
	Deploy map[string]string // Deploy information

	// Private fields. These are usually set on Get or Put.
	//
	// DO NOT MODIFY THESE.
	ID string
}

func (d *Deploy) setId() {
	d.ID = uuid.GenerateUUID()
}
