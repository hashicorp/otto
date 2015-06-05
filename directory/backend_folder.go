package directory

import (
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

func (b *FolderBackend) blobPath(k string) string {
	return filepath.Join(b.Dir, "blob", k)
}

func (b *FolderBackend) ensurePath(path string) error {
	return os.MkdirAll(path, 0755)
}
