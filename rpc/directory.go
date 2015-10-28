package rpc

import (
	"net/rpc"

	"github.com/hashicorp/otto/directory"
)

// Directory is an implementatin of directory.Backend that communicates
// over RPC.
type Directory struct {
	Client *rpc.Client
	Name   string
}

func (d *Directory) PutBlob(string, *directory.BlobData) error {
	return nil
}

func (d *Directory) GetBlob(string) (*directory.BlobData, error) {
	return nil, nil
}

func (d *Directory) PutInfra(v *directory.Infra) error {
	var resp ErrorResponse
	err := d.Client.Call(d.Name+".PutInfra", v, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

	return nil
}

func (d *Directory) GetInfra(v *directory.Infra) (*directory.Infra, error) {
	var resp DirGetInfraResponse
	err := d.Client.Call(d.Name+".GetInfra", v, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Value, nil
}

func (d *Directory) PutDev(v *directory.Dev) error {
	var resp ErrorResponse
	err := d.Client.Call(d.Name+".PutDev", v, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

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
	var resp ErrorResponse
	err := d.Client.Call(d.Name+".PutBuild", v, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

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
	var resp ErrorResponse
	err := d.Client.Call(d.Name+".PutDeploy", v, &resp)
	if err != nil {
		return err
	}
	if resp.Error != nil {
		err = resp.Error
		return err
	}

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

// DirectoryServer is a net/rpc compatible structure for serving
// a directory backend. This should not be used directly.
type DirectoryServer struct {
	Directory directory.Backend
}

func (s *DirectoryServer) PutInfra(
	args *directory.Infra,
	reply *ErrorResponse) error {
	err := s.Directory.PutInfra(args)
	*reply = ErrorResponse{
		Error: NewBasicError(err),
	}
	return nil
}

func (s *DirectoryServer) GetInfra(
	args *directory.Infra,
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
	reply *ErrorResponse) error {
	err := s.Directory.PutDev(args)
	*reply = ErrorResponse{
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
	reply *ErrorResponse) error {
	err := s.Directory.PutBuild(args)
	*reply = ErrorResponse{
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
	reply *ErrorResponse) error {
	err := s.Directory.PutDeploy(args)
	*reply = ErrorResponse{
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
