package bindata

import (
	"bytes"
	"io"
	"path/filepath"
)

// tplLoader is a pongo2.TemplateLoader that loads templates from assets.
type tplLoader struct {
	Data *Data
	Base string
}

func (t *tplLoader) Abs(base, name string) string {
	return filepath.Join(t.Base, name)
}

func (t *tplLoader) Get(path string) (io.Reader, error) {
	raw, err := t.Data.Asset(path)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(raw), nil
}
