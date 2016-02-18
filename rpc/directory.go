package rpc

import (
	"io"
	"log"
	"net/rpc"

	"github.com/hashicorp/otto/directory"
)

// Directory is an implementatin of directory.Backend that communicates
// over RPC.
type Directory struct {
	Broker *muxBroker
	Client *rpc.Client
	Name   string
}

func (d *Directory) PutApp(l *directory.AppLookup, v *directory.App) error {
	var resp DirPutAppResponse
	args := &DirPutAppArgs{Lookup: l, App: v}
	err := d.Client.Call(d.Name+".PutApp", args, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		return err
	}

	return nil
}

func (d *Directory) GetApp(l *directory.AppLookup) (*directory.App, error) {
	var resp DirGetAppResponse
	err := d.Client.Call(d.Name+".GetApp", l, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Value, nil
}

func (d *Directory) ListApps() ([]*directory.App, error) {
	var resp DirListAppsResponse
	err := d.Client.Call(d.Name+".ListApps", new(interface{}), &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Value, nil
}

func (d *Directory) PutInfra(l *directory.InfraLookup, infra *directory.Infra) error {
	var resp DirPutInfraResponse
	args := &DirPutInfraArgs{Lookup: l, Infra: infra}
	err := d.Client.Call(d.Name+".PutInfra", args, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

	return nil
}

func (d *Directory) GetInfra(l *directory.InfraLookup) (*directory.Infra, error) {
	var resp DirGetInfraResponse
	err := d.Client.Call(d.Name+".GetInfra", l, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Value, nil
}

func (d *Directory) PutBlob(key string, data *directory.BlobData) error {
	// Serve the data
	id := d.Broker.NextId()
	doneCh := make(chan struct{})
	go serveSingleCopy("putBlob: "+key, d.Broker, doneCh, id, nil, data.Data)

	// Run it
	var resp ErrorResponse
	args := &DirPutBlobArgs{
		Key: key,
		Id:  id,
	}
	err := d.Client.Call(d.Name+".PutBlob", args, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

	// NOTE: We don't wait on the doneCh because the data is being read
	// on the other side and the contract is that it will exaust the
	// data on success.

	return nil
}

func (d *Directory) GetBlob(key string) (*directory.BlobData, error) {
	// Create the result
	pr, pw := io.Pipe()
	result := &directory.BlobData{
		Key:  key,
		Data: pr,
	}

	// Download the data
	id := d.Broker.NextId()
	doneCh := make(chan struct{})
	go serveSingleCopy("getBlob: "+key, d.Broker, doneCh, id, pw, nil)
	go func() {
		// We wait for the data copying to be complete. When it is, we
		// want to close the writer side of our pipe so that the result
		// also gets an EOF.
		<-doneCh
		pw.Close()
	}()

	// Run it
	var resp DirGetBlobResponse
	args := &DirGetBlobArgs{
		Key: key,
		Id:  id,
	}
	err := d.Client.Call(d.Name+".GetBlob", args, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}
	if !resp.Ok {
		// No return value (blob key didn't exist)
		result = nil
	}

	// NOTE: We don't wait on the doneCh because the data is put as a
	// io.Reader into the BlobData. It is up to the end user to read
	// until it is exausted.

	return result, nil
}

func (d *Directory) PutDev(v *directory.Dev) error {
	var resp DirPutDevResponse
	err := d.Client.Call(d.Name+".PutDev", v, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

	*v = *resp.Value
	return nil
}

func (d *Directory) GetDev(v *directory.Dev) (*directory.Dev, error) {
	var resp DirGetDevResponse
	err := d.Client.Call(d.Name+".GetDev", v, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Value, nil
}

func (d *Directory) DeleteDev(v *directory.Dev) error {
	var resp ErrorResponse
	err := d.Client.Call(d.Name+".DeleteDev", v, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

	return nil
}

func (d *Directory) PutBuild(v *directory.Build) error {
	var resp DirPutBuildResponse
	err := d.Client.Call(d.Name+".PutBuild", v, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

	*v = *resp.Value
	return nil
}

func (d *Directory) GetBuild(v *directory.Build) (*directory.Build, error) {
	var resp DirGetBuildResponse
	err := d.Client.Call(d.Name+".GetBuild", v, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Value, nil
}

func (d *Directory) PutDeploy(v *directory.Deploy) error {
	var resp DirPutDeployResponse
	err := d.Client.Call(d.Name+".PutDeploy", v, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

	*v = *resp.Value
	return nil
}

func (d *Directory) GetDeploy(v *directory.Deploy) (*directory.Deploy, error) {
	var resp DirGetDeployResponse
	err := d.Client.Call(d.Name+".GetDeploy", v, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Value, nil
}

type DirPutBlobArgs struct {
	Key string
	Id  uint32
}

type DirGetBlobArgs struct {
	Key string
	Id  uint32
}

type DirPutAppArgs struct {
	Lookup *directory.AppLookup
	App    *directory.App
}

type DirGetBlobResponse struct {
	Ok    bool
	Error *BasicError
}

type DirGetAppResponse struct {
	Value *directory.App
	Error *BasicError
}

type DirListAppsResponse struct {
	Value []*directory.App
	Error *BasicError
}

type DirPutInfraArgs struct {
	Lookup *directory.InfraLookup
	Infra  *directory.Infra
}

type DirGetInfraResponse struct {
	Value *directory.Infra
	Error *BasicError
}

type DirGetDevResponse struct {
	Value *directory.Dev
	Error *BasicError
}

type DirGetBuildResponse struct {
	Value *directory.Build
	Error *BasicError
}

type DirGetDeployResponse struct {
	Value *directory.Deploy
	Error *BasicError
}

type DirPutAppResponse struct {
	Error *BasicError
}

type DirPutInfraResponse struct {
	Error *BasicError
}

type DirPutDevResponse struct {
	Value *directory.Dev
	Error *BasicError
}

type DirPutBuildResponse struct {
	Value *directory.Build
	Error *BasicError
}

type DirPutDeployResponse struct {
	Value *directory.Deploy
	Error *BasicError
}

// DirectoryServer is a net/rpc compatible structure for serving
// a directory backend. This should not be used directly.
type DirectoryServer struct {
	Broker    *muxBroker
	Directory directory.Backend
}

func (s *DirectoryServer) PutBlob(
	args *DirPutBlobArgs,
	reply *ErrorResponse) error {
	// Connect to the data stream
	conn, err := s.Broker.Dial(args.Id)
	if err != nil {
		*reply = ErrorResponse{Error: NewBasicError(err)}
		return nil
	}
	defer conn.Close()

	err = s.Directory.PutBlob(args.Key, &directory.BlobData{
		Data: conn,
	})
	*reply = ErrorResponse{
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) GetBlob(
	args *DirGetBlobArgs,
	reply *DirGetBlobResponse) error {
	// Connect to the data stream. We need to connect first because
	// we need to at least connect and close the connection otherwise
	// we'll leak a serveSingleCopy goroutine. This is why this cannot be
	// below the GetBlob call below.
	conn, err := s.Broker.Dial(args.Id)
	if err != nil {
		*reply = DirGetBlobResponse{Error: NewBasicError(err)}
		return nil
	}

	// Get the blob. If we have an error we return right away
	result, err := s.Directory.GetBlob(args.Key)
	*reply = DirGetBlobResponse{
		Ok:    result != nil,
		Error: NewBasicError(err),
	}
	if err != nil {
		conn.Close()
		return nil
	}
	if result == nil {
		// There is no data, so just close the data stream and return
		conn.Close()
		return nil
	}

	// No error! We need to copy the data from the blob into the
	// resulting connection. We do this in a goroutine though since the
	// data can be read at any speed.
	go func() {
		defer conn.Close()
		defer result.Close()
		if _, err := io.Copy(conn, result.Data); err != nil {
			log.Printf(
				"[ERR] rpc/directory: error copying getBlob data '%s': %s",
				args.Key,
				err)
		}
	}()

	return nil
}

func (s *DirectoryServer) PutApp(
	args *DirPutAppArgs,
	reply *DirPutAppResponse) error {
	err := s.Directory.PutApp(args.Lookup, args.App)
	*reply = DirPutAppResponse{
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) GetApp(
	args *directory.AppLookup,
	reply *DirGetAppResponse) error {
	result, err := s.Directory.GetApp(args)
	*reply = DirGetAppResponse{
		Value: result,
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) ListApps(
	args interface{},
	reply *DirListAppsResponse) error {
	result, err := s.Directory.ListApps()
	*reply = DirListAppsResponse{
		Value: result,
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) PutInfra(
	args *DirPutInfraArgs,
	reply *DirPutInfraResponse) error {
	err := s.Directory.PutInfra(args.Lookup, args.Infra)
	*reply = DirPutInfraResponse{
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) GetInfra(
	args *directory.InfraLookup,
	reply *DirGetInfraResponse) error {
	result, err := s.Directory.GetInfra(args)
	*reply = DirGetInfraResponse{
		Value: result,
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) PutDev(
	args *directory.Dev,
	reply *DirPutDevResponse) error {
	err := s.Directory.PutDev(args)
	*reply = DirPutDevResponse{
		Value: args,
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) GetDev(
	args *directory.Dev,
	reply *DirGetDevResponse) error {
	result, err := s.Directory.GetDev(args)
	*reply = DirGetDevResponse{
		Value: result,
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) DeleteDev(
	args *directory.Dev,
	reply *ErrorResponse) error {
	err := s.Directory.DeleteDev(args)
	*reply = ErrorResponse{
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) PutBuild(
	args *directory.Build,
	reply *DirPutBuildResponse) error {
	err := s.Directory.PutBuild(args)
	*reply = DirPutBuildResponse{
		Value: args,
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) GetBuild(
	args *directory.Build,
	reply *DirGetBuildResponse) error {
	result, err := s.Directory.GetBuild(args)
	*reply = DirGetBuildResponse{
		Value: result,
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) PutDeploy(
	args *directory.Deploy,
	reply *DirPutDeployResponse) error {
	err := s.Directory.PutDeploy(args)
	*reply = DirPutDeployResponse{
		Value: args,
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) GetDeploy(
	args *directory.Deploy,
	reply *DirGetDeployResponse) error {
	result, err := s.Directory.GetDeploy(args)
	*reply = DirGetDeployResponse{
		Value: result,
		Error: NewBasicError(err),
	}
	return nil
}
