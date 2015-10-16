package vagrant

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/otto/app"
)

// DevLayered returns a Layered setup for development.
//
// This automatically prepares any layers for foundations. Your custom
// layers are then appended to the foundation layers. If you have no layers,
// no modification is necessary to the returned value. It is ready to go
// as-is.
func DevLayered(ctx *app.Context, layers []*Layer) (*Layered, error) {
	// Basic result
	result := &Layered{
		DataDir: filepath.Join(ctx.GlobalDir, "vagrant-layered"),
		Layers:  make([]*Layer, 0, len(ctx.FoundationDirs)+len(layers)),
	}

	// Find all the foundation layers
	layersDir := filepath.Join(ctx.Dir, "foundation-layers")
	dir, err := os.Open(layersDir)
	if err != nil {
		// If we don't have foundation layers, we're done!
		if os.IsNotExist(err) {
			return result, nil
		}

		return nil, err
	}

	// Read the directory names
	dirs, err := dir.Readdirnames(-1)
	dir.Close()
	if err != nil {
		return nil, err
	}

	// Sort the directories so we get the right order
	sort.Strings(dirs)

	// Go through each directory and add the layer
	for _, dir := range dirs {
		parts := strings.SplitN(dir, "-", 2)
		id := parts[1]

		result.Layers = append(result.Layers, &Layer{
			ID:          id,
			Vagrantfile: filepath.Join(layersDir, dir, "Vagrantfile"),
		})
	}

	// Add our final layers
	if len(layers) > 0 {
		result.Layers = append(result.Layers, layers...)
	}

	return result, nil
}
