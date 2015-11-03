package rpc

import (
	"io"
	"net"
	"net/rpc"
	"testing"

	"github.com/hashicorp/otto/app"
)

func testConn(t *testing.T) (net.Conn, net.Conn) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	var serverConn net.Conn
	doneCh := make(chan struct{})
	go func() {
		defer close(doneCh)
		defer l.Close()
		var err error
		serverConn, err = l.Accept()
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	clientConn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	<-doneCh

	return clientConn, serverConn
}

func testClientServer(t *testing.T) (*rpc.Client, *rpc.Server) {
	clientConn, serverConn := testConn(t)

	server := rpc.NewServer()
	go server.ServeConn(serverConn)

	client := rpc.NewClient(clientConn)

	return client, server
}

func testNewClientServer(t *testing.T) (*Client, *Server, *testStreams) {
	clientConn, serverConn := testConn(t)

	server := &Server{
		AppFunc: testAppFixed(new(app.Mock)),
	}
	streams := testNewStreams(t, server)
	go server.ServeConn(serverConn)

	client, err := NewClient(clientConn)
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return client, server, streams
}

func testNewStreams(t *testing.T, s *Server) *testStreams {
	stdin_r, stdin_w := io.Pipe()
	stdout_r, stdout_w := io.Pipe()
	stderr_r, stderr_w := io.Pipe()

	s.Stdin = stdin_w
	s.Stdout = stdout_r
	s.Stderr = stderr_r

	return &testStreams{
		Stdin:  stdin_r,
		Stdout: stdout_w,
		Stderr: stderr_w,
	}
}

func testAppFixed(c app.App) AppFunc {
	return func() app.App {
		return c
	}
}

type testStreams struct {
	Stdin  io.ReadCloser
	Stdout io.WriteCloser
	Stderr io.WriteCloser
}

func (s *testStreams) Close() error {
	s.Stdin.Close()
	s.Stdout.Close()
	s.Stderr.Close()
	return nil
}
