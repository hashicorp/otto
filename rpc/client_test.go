package rpc

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/hashicorp/otto/app"
)

func TestClient_App(t *testing.T) {
	clientConn, serverConn := testConn(t)

	c := new(app.Mock)
	server := &Server{AppFunc: testAppFixed(c)}
	streams := testNewStreams(t, server)
	defer streams.Close()

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

func TestClient_syncStreams(t *testing.T) {
	client, _, streams := testNewClientServer(t)

	// Start the data copying
	var stdout_out, stderr_out, stdin_out bytes.Buffer
	stdout := bytes.NewBufferString("stdouttest")
	stderr := bytes.NewBufferString("stderrtest")
	stdin := bytes.NewBufferString("stdintest")
	go client.SyncStreams(stdin, &stdout_out, &stderr_out)
	go io.Copy(&stdin_out, streams.Stdin)
	go io.Copy(streams.Stdout, stdout)
	go io.Copy(streams.Stderr, stderr)

	// Unfortunately I can't think of a better way to make sure all the
	// copies above go through so let's just exit.
	time.Sleep(100 * time.Millisecond)

	// Close everything, and lets test the result
	client.Close()
	streams.Close()

	if v := stdin_out.String(); v != "stdintest" {
		t.Fatalf("bad: %s", v)
	}
	if v := stdout_out.String(); v != "stdouttest" {
		t.Fatalf("bad: %s", v)
	}
	if v := stderr_out.String(); v != "stderrtest" {
		t.Fatalf("bad: %s", v)
	}
}
