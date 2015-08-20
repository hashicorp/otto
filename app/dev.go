package app

// DevDep has information about an upstream dependency that should be
// used by the Dev function in order to build a complete development
// environment.
type DevDep struct {
	// FragmentPath is the path to the file that contains the
	// Vagrantfile fragment necessary to configure and run this
	// dependency.
	FragmentPath string
}
