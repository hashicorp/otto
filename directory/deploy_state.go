package directory

//go:generate stringer -type=DeployState deploy_state.go

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
