package directory

// Infra represents the data stored in the directory service about
// Infrastructures.
type Infra struct {
	InfraLookup // InfraLookup is the lookup data for this Infra

	Name   string // Name of this infra
	Type   string // Type of this infra
	Flavor string // Flavor of the infra

	// DeployVersion is the deploy version of this infra. This MUST
	// be a semantic version. If it doesn't parse as a semantic version,
	// errors will be raised.
	DeployVersion string

	// State is the state of this infrastructure. This is important since
	// it is possible for there to be a partial state. If we're in a
	// partial state then deploys and such can't go through yet.
	State InfraState `json:"state"`

	// Outputs are the output data from the infrastructure step.
	// This is an opaque blob that is dependent on each infrastructure
	// type. Please refer to docs of a specific infra to learn more about
	// what values are here.
	Outputs map[string]string `json:"outputs"`

	// Opaque is extra data associated with this infrastructure. Anything
	// can be stored here but it should be minimal, if possible since the
	// directory backends aren't meant to be large binary storage.
	Opaque []byte
}

// InfraLookup is the structure used to look up or store infras.
type InfraLookup struct {
	Name string // Name of the infrastructure
}
