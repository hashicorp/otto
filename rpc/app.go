package rpc

import (
	"net/rpc"

	"github.com/hashicorp/otto/app"
)

// App is an implementation of app.App that communicates over RPC.
type App struct {
	Broker *muxBroker
	Client *rpc.Client
	Name   string
}

func (c *App) Compile(ctx *app.Context) (*app.CompileResult, error) {
	var resp AppCompileResponse
	args := AppCompileArgs{Context: ctx}

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

func (c *App) Build(ctx *app.Context) error                      { return nil }
func (c *App) Deploy(ctx *app.Context) error                     { return nil }
func (c *App) Dev(ctx *app.Context) error                        { return nil }
func (c *App) DevDep(dst, src *app.Context) (*app.DevDep, error) { return nil, nil }

func (c *App) Close() error {
	return c.Client.Close()
}

// AppServer is a net/rpc compatible structure for serving an App.
// This should not be used directly.
type AppServer struct {
	Broker *muxBroker
	App    app.App
}

type AppCompileArgs struct {
	ContextSharedArgs

	Context *app.Context
}

type AppCompileResponse struct {
	Result *app.CompileResult
	Error  *BasicError
}

func (s *AppServer) Compile(
	args *AppCompileArgs,
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
