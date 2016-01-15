package directory

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/terraform/dag"
)

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

// AppSlice is a wrapper around []*App that implements sort.Interface.
// The sorting order is standard sorting for the tuple:
// (app name, app ID, version)
type AppSlice []*App

func (a AppSlice) Len() int {
	return len(a)
}

func (a AppSlice) Less(i, j int) bool {
	// Order: (name, ID, version)

	if a[i].Name != a[j].Name {
		return a[i].Name < a[j].Name
	}

	if a[i].AppLookup.AppID != a[j].AppLookup.AppID {
		return a[i].AppLookup.AppID < a[j].AppLookup.AppID
	}

	if a[i].AppLookup.Version != a[j].AppLookup.Version {
		// Parse the versions. We panic if these fail since we should
		// verify prior to inserting into the directory that these are valid.
		v1, e1 := version.NewVersion(a[i].AppLookup.Version)
		v2, e2 := version.NewVersion(a[j].AppLookup.Version)
		if e1 != nil {
			panic(e1)
		}
		if e2 != nil {
			panic(e2)
		}

		return v1.LessThan(v2)
	}

	// Equal
	return false
}

func (a AppSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
