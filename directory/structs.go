package directory

// Environment is the structure used to represent an environment where
// applications are deployed, such as "staging".
type Environment struct {
	Name  string      // Name is the name of the environment, i.e. "staging"
	Infra InfraLookup // Infra is the infra that this environment uses
}

// InfraLookup is the structure used to look up or store an infra.
type InfraLookup struct {
	Name   string // Name is the name of the infra
	Type   string // Type is the type of infra, i.e. "aws"
	Flavor string // Flavor is the flavor of the infra
}
