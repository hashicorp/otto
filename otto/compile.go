package otto

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/foundation"
	"github.com/hashicorp/otto/infrastructure"
)

// CompileMetadata is the stored metadata about a successful compilation.
//
// Failures during compilation result in no metadata at all being stored.
// This metadata can be used to access various information about the resulting
// compilation.
type CompileMetadata struct {
	// App is the result of compiling the main application
	App *app.CompileResult `json:"app"`

	// Deps are the results of compiling the dependencies, keyed by their
	// unique Otto ID. If you want the tree structure then use the Appfile
	// itself to search the dependency tree, then the ID of that dep
	// to key into this map.
	AppDeps map[string]*app.CompileResult `json:"app_deps"`

	// Infra is the result of compiling the infrastructure for this application
	Infra *infrastructure.CompileResult `json:"infra"`

	// Foundations is the listing of top-level foundation compilation results.
	Foundations map[string]*foundation.CompileResult `json:"foundations"`
}

func (c *Core) resetCompileMetadata() {
	c.metadataCache = nil
}

func (c *Core) compileMetadata() (*CompileMetadata, error) {
	if c.metadataCache != nil {
		return c.metadataCache, nil
	}

	f, err := os.Open(filepath.Join(c.compileDir, "metadata.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}
	defer f.Close()

	var result CompileMetadata
	dec := json.NewDecoder(f)
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	c.metadataCache = &result
	return &result, nil
}

func (c *Core) saveCompileMetadata(md *CompileMetadata) error {
	if err := os.MkdirAll(c.compileDir, 0755); err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(c.compileDir, "metadata.json"))
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(md)
}
