// Package scriptpack is used to work with "ScriptPacks" which are
// packages of shell scripts that app types can define and use to get things
// done on the remote machine, whether its for development or deployment.
//
// ScriptPacks are 100% pure shell scripting. Any inputs must be received from
// environment variables. They aren't allowed to template at all. This is
// all done to ensure testability of the ScriptPacks.
//
// These are treated as first class elements within Otto to assist with
// testing.
//
// To create your own scriptpack, see the "template" folder within this
// directory. The folder structure and contents are important for scriptpacks
// to function correctly.
package scriptpack

// ScriptPack is a struct representing a single ScriptPack. This is exported
// from the various scriptpacks.
type ScriptPack struct {
	// Name is an identifying name used to name the environment variables.
	// For example the root of where the files are will always be
	// SCRIPTPACK_<NAME>_ROOT.
	Name string

	// Data is the compiled bindata for this ScriptPack. The entire
	// AssetDirectory will be copied starting from "/data"
	Data bindata.Data

	// Dependencies are a list of other ScriptPacks that will always be
	// unpacked alongside this ScriptPack. By referencing actual ScriptPack
	// pointers, the dependencies will also be statically compiled into
	// the Go binaries that contain them.
	//
	// Dependencies can be accessed at the path specified by
	// SCRIPTPACK_<DEP>_ROOT. Note that if you depend on a ScriptPack
	// which itself has a conflicting named dependency, then the first
	// one loaded will win. Be careful about this.
	Dependencies []*ScriptPack
}
