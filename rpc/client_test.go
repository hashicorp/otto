package rpc

import (
	"testing"

	"github.com/hashicorp/otto/app"
)

func TestClient_App(t *testing.T) {
	clientConn, serverConn := testConn(t)

	c := new(app.Mock)
	server := &Server{AppFunc: testAppFixed(c)}
	go server.ServeConn(serverConn)

	client, err := NewClient(clientConn)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer client.Close()

	appReal, err := client.App()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	// Compile!
	_, err = appReal.Compile(new(app.Context))
	if !c.CompileCalled {
		t.Fatal("compile should be called")
	}
	if err != nil {
		t.Fatalf("bad: %#v", err)
	}
}
