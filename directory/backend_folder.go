package directory

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// FolderBackend is a Directory backend that stores data on local disk.
//
// The primary use case for the FolderBackend is out-of-box experience
// for Otto and single developers. For team usage, FolderBackend is
// not recommended.
type FolderBackend struct {
	// Directory where data will be written. This directory will be
	// created if it doesn't exist.
	Dir string
}

func (b *FolderBackend) GetBlob(k string) (*BlobData, error) {
	path := b.blobPath(k)
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &BlobData{
		Key:    k,
		Data:   f,
		closer: f,
	}, nil
}

func (b *FolderBackend) PutBlob(k string, d *BlobData) error {
	path := b.blobPath(k)
	if err := b.ensurePath(filepath.Dir(path)); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, d.Data)
	return err
}

func (b *FolderBackend) GetInfra(k string) (*Infra, error) {
	var infra Infra
	ok, err := b.getData(b.infraPath(k), &infra)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &infra, nil
}

func (b *FolderBackend) PutInfra(k string, infra *Infra) error {
	return b.putData(b.infraPath(k), infra)
}

func (b *FolderBackend) blobPath(k string) string {
	return filepath.Join(b.Dir, "blob", k)
}

func (b *FolderBackend) infraPath(k string) string {
	return filepath.Join(b.Dir, "infra", k)
}

func (b *FolderBackend) getData(path string, d interface{}) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}

		return false, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	return true, dec.Decode(d)
}

func (b *FolderBackend) putData(path string, d interface{}) error {
	if err := b.ensurePath(filepath.Dir(path)); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Because this data is meant to be debuggable, let's output
	// it in human-readable format.
	data, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}

	// Copy the data directly into the file
	_, err = io.Copy(f, bytes.NewReader(data))
	return err
}

func (b *FolderBackend) ensurePath(path string) error {
	return os.MkdirAll(path, 0755)
}
