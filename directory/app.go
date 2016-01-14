package directory

import (
	"fmt"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/terraform/dag"
)

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

// AppLookupAppfile creates an AppLookup from an Appfile
func AppLookupAppfile(f *appfile.File) *AppLookup {
	return &AppLookup{
		AppID:      f.ID,
		Version:    f.Application.Version().String(),
		ConfigHash: f.ConfigHash(),
	}
}

// App represents the data stored in the directory for a single
// application (a single Appfile).
type App struct {
	AppLookup // AppLookup is the lookup data for this App.

	Name         string      // Name of this application
	Type         string      // Type of this application
	Dependencies []AppLookup // Dependencies this app depends on
}

// NewAppCompiled creates an App instance from a compiled Appfile with
// the given root. The root must exist in the graph or it is an error.
func NewAppCompiled(c *appfile.Compiled, root dag.Vertex) (*App, error) {
	// Verify the vertex is in this graph
	if !c.Graph.HasVertex(root) {
		return nil, fmt.Errorf("%s not in Appfile graph", dag.VertexName(root))
	}

	// Get the typed graph vertex so we have access to the value
	gv, ok := root.(*appfile.CompiledGraphVertex)
	if !ok {
		panic(fmt.Sprintf("unknown graph vertex type: %T", root))
	}

	// TODO: dependencies

	return &App{
		AppLookup: *AppLookupAppfile(gv.File),
		Name:      gv.File.Application.Name,
		Type:      gv.File.Application.Type,
	}, nil
}
