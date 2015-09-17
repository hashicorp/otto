package directory

//go:generate stringer -type=InfraState infra_state.go

// InfraState is used to track the state of an infrastructure.
//
// This is required because the state of an infrastructure isn't binary.
// It can be not created at all, partially created, or fully created. We
// need to represent this in the directory.
type InfraState byte

const (
	InfraStateInvalid InfraState = 0
	InfraStatePartial InfraState = iota
	InfraStateReady
)
