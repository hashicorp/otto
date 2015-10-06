package bindata

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
)

// tplLoader is a pongo2.TemplateLoader that loads templates from assets.
type tplLoader struct {
	Data *Data
	Base string
}

func (t *tplLoader) Abs(base, name string) string {
	idx := strings.IndexAny(name, ":/")
	if idx != -1 && name[idx] == ':' {
		// We have a shared resource. Return the name as-is
		return name
	}

	return filepath.Join(t.Base, name)
}

func (t *tplLoader) Get(path string) (io.Reader, error) {
	data := t.Data

	idx := strings.IndexAny(path, ":/")
	if idx != -1 && path[idx] == ':' {
		// We have a shared resource. Get the proper data loader.
		share := path[:idx]
		raw, ok := t.Data.SharedExtends[share]
		if !ok {
			return nil, fmt.Errorf("share not found: %s", share)
		}

		data = raw
		path = path[idx+1:]
	}

	raw, err := data.Asset(path)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(raw), nil
}
