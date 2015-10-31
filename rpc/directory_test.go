package rpc

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/directory"
)

func TestDirectory_impl(t *testing.T) {
	var _ directory.Backend = new(Directory)
}

func TestDirectory(t *testing.T) {
	// Create the temporary directory for the directory data
	td, err := ioutil.TempDir("", "otto")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer os.RemoveAll(td)

	// Create the actual plugin client/server
	client, server := testNewClientServer(t)
	defer client.Close()

	// Build a context that points to our bolt directory backend
	ctx := new(app.Context)
	ctx.Directory = &directory.BoltBackend{Dir: td}

	// Create an appMock. The mock has a compile function that actually
	// calls the directory test on it to verify that this works properly.
	//
	// This will verify the backend that is being passed through the
	// RPC layer actually works. We have to test it within the callback
	// since the connection is over after that point.
	appMock := server.AppFunc().(*app.Mock)
	appReal, err := client.App()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	appMock.CompileFunc = func(ctx *app.Context) (r *app.CompileResult, err error) {
		directory.TestBackend(t, ctx.Directory)
		return
	}

	_, err = appReal.Compile(ctx)
	if !appMock.CompileCalled {
		t.Fatal("compile should be called")
	}
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
}
