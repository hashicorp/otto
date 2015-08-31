package directory

import (
	"io"
	"os"
)

// Backend is the interface for any directory service. It is effectively
// the protocol right now. In the future we may abstract this further to
// an official protocol if it proves to be necessary. Until then, it is
// a value add on top of the Appfile (but not part of that format) that Otto
// uses for global state.
type Backend interface {
	// PutBlob writes binary data for a given project/infra/app.
	//
	// GetBlob reads that data back out.
	//
	// ListBlob lists the binary data stored.
	PutBlob(string, *BlobData) error
	GetBlob(string) (*BlobData, error)

	// PutInfra and GetInfra are the functions used to store and retrieve
	// data about infrastructures.
	PutInfra(*Infra) error
	GetInfra(*Infra) (*Infra, error)

	// PutBuild stores the result of a build.
	//
	// GetBuild queries a build. The result is returned. The parameter
	// must fill in the App, Infra, and InfraFlavor fields.
	PutBuild(*Build) error
	GetBuild(*Build) (*Build, error)

	// PutDeploy stores the result of a build.
	//
	// GetDeploy queries a deploy. The result is returned. The parameter
	// must fill in the App, Infra, and InfraFlavor fields.
	PutDeploy(*Deploy) error
	GetDeploy(*Deploy) (*Deploy, error)
}

// Build represents a build of an App.
type Build struct {
	App         string            // App is the app type, i.e. "go"
	Infra       string            // Infra is the infra type, i.e. "aws"
	InfraFlavor string            // InfraFlavor is the flavor, i.e. "vpc-public-private"
	Artifact    map[string]string // Resulting artifact from the build
}

// BlobData is the metadata and data associated with stored binary
// data. The fields and their usage varies depending on the operations,
// so please read the documentation for each field carefully.
type BlobData struct {
	// Key is the key for the blob data. This is populated on read and ignored
	// on any other operation.
	Key string

	// Data is the data for the blob data. When writing, this should be
	// the data to write. When reading, this is the data that is read.
	Data io.Reader

	closer io.Closer
}

func (d *BlobData) Close() error {
	if d.closer != nil {
		return d.closer.Close()
	}

	return nil
}

// WriteToFile is a helper to write BlobData to a file. While this is
// a very easy thing to do, it is so common that we provide a function
// for doing so.
func (d *BlobData) WriteToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, d.Data)
	return err
}
