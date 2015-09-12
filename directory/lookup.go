package directory

// Lookup has fields that are used for looking up data in the directory.
//
// Note that not all fields are required for every lookup operation. Please
// see the documentation for the item you're trying to look up to learn more.
type Lookup struct {
	App         string // App is the app type, i.e. "go"
	Infra       string // Infra is the infra type, i.e. "aws"
	InfraFlavor string // InfraFlavor is the flavor, i.e. "vpc-public-private"
	Foundation  string // Foundation is the name of he foundation, i.e. "consul"
}
