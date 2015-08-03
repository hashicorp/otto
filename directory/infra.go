package directory

// Infra represents the data stored in the directory service about
// Infrastructures.
type Infra struct {
	// Outputs are the output data from the infrastructure step.
	// This is an opaque blob that is dependent on each infrastructure
	// type. Please refer to docs of a specific infra to learn more about
	// what values are here.
	Outputs map[string]string `json:"outputs"`
}
