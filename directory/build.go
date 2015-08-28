package directory

// Build represents a build of an App.
type Build struct {
	App         string            // App is the app type, i.e. "go"
	Infra       string            // Infra is the infra type, i.e. "aws"
	InfraFlavor string            // InfraFlavor is the flavor, i.e. "vpc-public-private"
	Artifact    map[string]string // Resulting artifact from the build
}
