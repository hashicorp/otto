package directory

//go:generate stringer -type=DevState infra_state.go

// DevState is used to track the state of an infrastructure.
//
// This is required because the state of an infrastructure isn't binary.
// It can be not created at all, partially created, or fully created. We
// need to represent this in the directory.
type DevState byte

const (
	DevStateInvalid DevState = 0
	DevStateNew     DevState = iota
	DevStateReady
)
