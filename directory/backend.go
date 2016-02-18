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
	//--------------------------------------------------------------------
	// App-related
	//--------------------------------------------------------------------

	// PutApp stores the application in the backend. This denotes that
	// an application exists but may not yet be deployed, and codifies
	// the configuration that is in use. If the application already exists,
	// this will overwrite it.
	//
	// GetApp finds an exact application using some lookup values. If the
	// app is not found, nil is returned.
	//
	// ListApps returns the list of Apps that are available in the directory.
	// These will be sorted according to AppSlice.
	PutApp(*AppLookup, *App) error
	GetApp(*AppLookup) (*App, error)
	ListApps() ([]*App, error)

	//--------------------------------------------------------------------
	// Infra-related
	//--------------------------------------------------------------------

	// PutInfra stores the state of an infrastructure.
	//
	// GetInfra finds an exact Infra matching the lookup values. If the
	// infra is not found, nil is returned.
	//
	// ListInfra returns the list of infrastructures that are availabile
	// in the directory. This will be sorted according to InfraSlice.
	PutInfra(*InfraLookup, *Infra) error
	GetInfra(*InfraLookup) (*Infra, error)
	ListInfra() ([]*Infra, error)

	//--------------------------------------------------------------------
	// Legacy
	//--------------------------------------------------------------------

	// PutBlob writes binary data for a given project/infra/app.
	//
	// GetBlob reads that data back out.
	//
	// ListBlob lists the binary data stored.
	PutBlob(string, *BlobData) error
	GetBlob(string) (*BlobData, error)

	// PutDev stores the result of a dev.
	//
	// GetDev queries a dev. The result is returned. The parameter
	// must fill in the App, Infra, and InfraFlavor fields.
	PutDev(*Dev) error
	GetDev(*Dev) (*Dev, error)
	DeleteDev(*Dev) error

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
	// Lookup information for the Build. AppID, Infra, and InfraFlavor
	// are required.
	Lookup

	// Resulting artifact from the build
	Artifact map[string]string
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
