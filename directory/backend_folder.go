package directory

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (b *FolderBackend) GetInfra(infra *Infra) (*Infra, error) {
	var result Infra
	ok, err := b.getData(b.infraPath(infra), &result)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &result, nil
}

func (b *FolderBackend) PutInfra(infra *Infra) error {
	if infra.ID == "" {
		infra.setId()
	}

	return b.putData(b.infraPath(infra), infra)
}

func (b *FolderBackend) PutBuild(build *Build) error {
	return b.putData(b.buildPath(build), build)
}

func (b *FolderBackend) GetBuild(build *Build) (*Build, error) {
	var result Build
	ok, err := b.getData(b.buildPath(build), &result)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &result, nil
}

func (b *FolderBackend) PutDeploy(deploy *Deploy) error {
	if deploy.ID == "" {
		deploy.setId()
	}

	return b.putData(b.deployPath(deploy), deploy)
}

func (b *FolderBackend) GetDeploy(deploy *Deploy) (*Deploy, error) {
	var result Deploy
	ok, err := b.getData(b.deployPath(deploy), &result)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &result, nil
}

func (b *FolderBackend) blobPath(k string) string {
	return filepath.Join(b.Dir, "blob", k)
}

func (b *FolderBackend) buildPath(build *Build) string {
	return filepath.Join(
		b.Dir,
		"build",
		fmt.Sprintf("%s-%s-%s", build.App, build.Infra, build.InfraFlavor))
}

func (b *FolderBackend) deployPath(deploy *Deploy) string {
	return filepath.Join(
		b.Dir,
		"deploy",
		fmt.Sprintf("%s-%s-%s", deploy.App, deploy.Infra, deploy.InfraFlavor))
}

func (b *FolderBackend) infraPath(infra *Infra) string {
	key := infra.Lookup.Infra
	if infra.Lookup.Foundation != "" {
		key += "-" + infra.Lookup.Foundation
	}

	return filepath.Join(b.Dir, "infra", key)
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
