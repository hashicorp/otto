package rpc

import (
	"reflect"
	"testing"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/otto/ui"
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

func TestApp_compileUi(t *testing.T) {
	client, server := testNewClientServer(t)
	defer client.Close()

	appMock := server.AppFunc().(*app.Mock)
	appReal, err := client.App()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	appMock.CompileFunc = func(ctx *app.Context) (*app.CompileResult, error) {
		ctx.Ui.Message("HELLO!")
		return nil, nil
	}

	ui := new(ui.Mock)
	ctx := new(app.Context)
	ctx.Ui = ui

	_, err = appReal.Compile(ctx)
	if !appMock.CompileCalled {
		t.Fatal("compile should be called")
	}
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}

	if ui.MessageBuf[0] != "HELLO!" {
		t.Fatalf("bad: %#v", ui)
	}
}
