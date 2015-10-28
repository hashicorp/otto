package rpc

import (
	"reflect"
	"testing"

	"github.com/hashicorp/otto/app"
)

func TestApp_impl(t *testing.T) {
	var _ app.App = new(App)
}

func TestApp_compile(t *testing.T) {
	client, server := testNewClientServer(t)
	defer client.Close()

	appMock := server.AppFunc().(*app.Mock)
	appReal, err := client.App()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	appMock.CompileResult = &app.CompileResult{Version: 42}

	actual, err := appReal.Compile(new(app.Context))
	if !appMock.CompileCalled {
		t.Fatal("compile should be called")
	}
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}

	expected := appMock.CompileResult
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("bad: %#v", actual)
	}
}
