package directory

// AppLookup is the structure used to look up or store an application.
//
// Some fields are ignored/unused for certain operations. See the documentation
// for the function using the structure for information.
type AppLookup struct {
	// Unique identifying fields: used for specific CRUD
	AppID      string // Otto-generated app UUID
	Version    string // Current version
	ConfigHash uint64 // Unique hash of the configuration, see Appfile.ConfigHash

	// Search fields: used for searching
	VersionConstraint string // Lookup based on constraints
}

// App represents the data stored in the directory for a single
// application (a single Appfile).
type App struct {
	AppLookup // AppLookup is the lookup data for this App.

	Name         string      // Name of this application
	Type         string      // Type of this application
	Dependencies []AppLookup // Dependencies this app depends on
}

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
