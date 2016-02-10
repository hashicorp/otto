package otto

import (
	"path/filepath"

	"github.com/hashicorp/otto/appfile"
	"github.com/hashicorp/otto/context"
)

//--------------------------------------------------------------------
// Core Methods
//--------------------------------------------------------------------

// contextShared returns the shared context structure for the given Appfile.
//
// Note that the shared context data may not be completely populated.
//
// If f is nil, the root Appfile is automatically used.
func (c *Core) contextShared(f *appfile.File) (*context.Shared, error) {
	// If f is nil, get the root Appfile
	if f == nil {
		root, err := c.appfileCompiled.Graph.Root()
		if err != nil {
			return nil, err
		}

		f = root.(*appfile.CompiledGraphVertex).File
	}

	return &context.Shared{
		Appfile:    f,
		InstallDir: filepath.Join(c.dataDir, "binaries"),
		Directory:  c.dir,
		Ui:         c.ui,
	}, nil
}
