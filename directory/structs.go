package directory

// Environment is the structure used to represent an environment where
// applications are deployed, such as "staging".
type Environment struct {
	Name  string      // Name is the name of the environment, i.e. "staging"
	Infra InfraLookup // Infra is the infra that this environment uses
}
