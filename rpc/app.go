package rpc

import (
	"net/rpc"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/appfile"
)

// App is an implementation of app.App that communicates over RPC.
type App struct {
	Broker *muxBroker
	Client *rpc.Client
	Name   string
}

func (c *App) Meta() (*app.Meta, error) {
	var resp AppMetaResponse

	err := c.Client.Call(c.Name+".Meta", new(struct{}), &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Result, nil
}

func (c *App) Implicit(ctx *app.Context) (*appfile.File, error) {
	var resp AppImplicitResponse
	args := AppContextArgs{Context: ctx}

	// Serve the shared context data
	serveContext(c.Broker, &ctx.Shared, &args.ContextSharedArgs)

	// Call
	err := c.Client.Call(c.Name+".Implicit", &args, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Result, nil
}

func (c *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	var resp AppCompileResponse
	args := AppContextArgs{Context: ctx}

	// Serve the shared context data
	serveContext(c.Broker, &ctx.Shared, &args.ContextSharedArgs)

	// Call
	err := c.Client.Call(c.Name+".Compile", &args, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		err = resp.Error
		return nil, err
	}

	return resp.Result, nil
}

func (c *App) Build(ctx *app.Context) error {
	var resp AppSimpleResponse
	args := AppContextArgs{Context: ctx}

	// Serve the shared context data
	serveContext(c.Broker, &ctx.Shared, &args.ContextSharedArgs)

	// Call
	err := c.Client.Call(c.Name+".Build", &args, &resp)
	if err == nil {
		if resp.Error != nil {
			err = resp.Error
		}
	}

	return err
}

func (c *App) Deploy(ctx *app.Context) error {
	var resp AppSimpleResponse
	args := AppContextArgs{Context: ctx}

	// Serve the shared context data
	serveContext(c.Broker, &ctx.Shared, &args.ContextSharedArgs)

	// Call
	err := c.Client.Call(c.Name+".Deploy", &args, &resp)
	if err == nil {
		if resp.Error != nil {
			err = resp.Error
		}
	}

	return err
}

func (c *App) Dev(ctx *app.Context) error {
	var resp AppSimpleResponse
	args := AppContextArgs{Context: ctx}

	// Serve the shared context data
	serveContext(c.Broker, &ctx.Shared, &args.ContextSharedArgs)

	// Call
	err := c.Client.Call(c.Name+".Dev", &args, &resp)
	if err == nil {
		if resp.Error != nil {
			err = resp.Error
		}
	}

	return err
}

func (c *App) DevDep(dst, src *app.Context) (*app.DevDep, error) {
	var resp AppDevDepResponse
	args := AppDevDepArgs{
		ContextDst: dst,
		ContextSrc: src,
	}

	// Serve the shared context data
	serveContext(c.Broker, &dst.Shared, &args.ContextDstShared)
	serveContext(c.Broker, &src.Shared, &args.ContextSrcShared)

	// Call
	err := c.Client.Call(c.Name+".DevDep", &args, &resp)
	if err == nil && resp.Error != nil {
		err = resp.Error
	}
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func (c *App) Close() error {
	return c.Client.Close()
}

// AppServer is a net/rpc compatible structure for serving an App.
// This should not be used directly.
type AppServer struct {
	Broker *muxBroker
	App    app.App
}

type AppContextArgs struct {
	ContextSharedArgs

	Context *app.Context
}

type AppMetaResponse struct {
	Result *app.Meta
	Error  *BasicError
}

type AppCompileResponse struct {
	Result *app.CompileResult
	Error  *BasicError
}

type AppImplicitResponse struct {
	Result *appfile.File
	Error  *BasicError
}

type AppDevDepArgs struct {
	ContextDst       *app.Context
	ContextDstShared ContextSharedArgs
	ContextSrc       *app.Context
	ContextSrcShared ContextSharedArgs
}

type AppDevDepResponse struct {
	Result *app.DevDep
	Error  *BasicError
}

type AppSimpleResponse struct {
	Error *BasicError
}

func (s *AppServer) Meta(
	args *struct{},
	reply *AppMetaResponse) error {
	result, err := s.App.Meta()
	*reply = AppMetaResponse{
		Result: result,
		Error:  NewBasicError(err),
	}

	return nil
}

func (s *AppServer) Implicit(
	args *AppContextArgs,
	reply *AppImplicitResponse) error {
	closer, err := connectContext(s.Broker, &args.Context.Shared, &args.ContextSharedArgs)
	defer closer.Close()
	if err != nil {
		*reply = AppImplicitResponse{
			Error: NewBasicError(err),
		}

		return nil
	}

	result, err := s.App.Implicit(args.Context)
	*reply = AppImplicitResponse{
		Result: result,
		Error:  NewBasicError(err),
	}

	return nil
}

func (s *AppServer) Compile(
	args *AppContextArgs,
	reply *AppCompileResponse) error {
	closer, err := connectContext(s.Broker, &args.Context.Shared, &args.ContextSharedArgs)
	defer closer.Close()
	if err != nil {
		*reply = AppCompileResponse{
			Error: NewBasicError(err),
		}

		return nil
	}

	result, err := s.App.Compile(args.Context)
	*reply = AppCompileResponse{
		Result: result,
		Error:  NewBasicError(err),
	}

	return nil
}

func (s *AppServer) Build(
	args *AppContextArgs,
	reply *AppSimpleResponse) error {
	closer, err := connectContext(s.Broker, &args.Context.Shared, &args.ContextSharedArgs)
	defer closer.Close()
	if err != nil {
		*reply = AppSimpleResponse{
			Error: NewBasicError(err),
		}

		return nil
	}

	*reply = AppSimpleResponse{
		Error: NewBasicError(s.App.Build(args.Context)),
	}

	return nil
}

func (s *AppServer) Deploy(
	args *AppContextArgs,
	reply *AppSimpleResponse) error {
	closer, err := connectContext(s.Broker, &args.Context.Shared, &args.ContextSharedArgs)
	defer closer.Close()
	if err != nil {
		*reply = AppSimpleResponse{
			Error: NewBasicError(err),
		}

		return nil
	}

	*reply = AppSimpleResponse{
		Error: NewBasicError(s.App.Deploy(args.Context)),
	}

	return nil
}

func (s *AppServer) Dev(
	args *AppContextArgs,
	reply *AppSimpleResponse) error {
	closer, err := connectContext(s.Broker, &args.Context.Shared, &args.ContextSharedArgs)
	defer closer.Close()
	if err != nil {
		*reply = AppSimpleResponse{
			Error: NewBasicError(err),
		}

		return nil
	}

	*reply = AppSimpleResponse{
		Error: NewBasicError(s.App.Dev(args.Context)),
	}

	return nil
}

func (s *AppServer) DevDep(
	args *AppDevDepArgs,
	reply *AppDevDepResponse) error {
	// Dst
	closer, err := connectContext(s.Broker, &args.ContextDst.Shared, &args.ContextDstShared)
	defer closer.Close()
	if err != nil {
		*reply = AppDevDepResponse{
			Error: NewBasicError(err),
		}

		return nil
	}

	// Src
	closer, err = connectContext(s.Broker, &args.ContextSrc.Shared, &args.ContextSrcShared)
	defer closer.Close()
	if err != nil {
		*reply = AppDevDepResponse{
			Error: NewBasicError(err),
		}

		return nil
	}

	result, err := s.App.DevDep(args.ContextDst, args.ContextSrc)
	*reply = AppDevDepResponse{
		Result: result,
		Error:  NewBasicError(err),
	}

	return nil
}
