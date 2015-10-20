package otto

import (
	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/infrastructure"
)

// CompileMetadata is the stored metadata about a successful compilation.
//
// Failures during compilation result in no metadata at all being stored.
// This metadata can be used to access various information about the resulting
// compilation.
type CompileMetadata struct {
	// App is the result of compiling the main application
	App *app.CompileResult

	// Deps are the results of compiling the dependencies, keyed by their
	// unique Otto ID. If you want the tree structure then use the Appfile
	// itself to search the dependency tree, then the ID of that dep
	// to key into this map.
	AppDeps map[string]*app.CompileResult

	// Infra is the result of compiling the infrastructure for this application
	Infra *infrastructure.CompileResult
}
