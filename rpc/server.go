package rpc

import (
	"io"
	"log"
	"net"
	"net/rpc"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/yamux"
)

// Server listens for network connections and then dispenses interface
// implementations over net/rpc.
type Server struct {
	AppFunc AppFunc

	// Stdout, Stderr are what this server will use instead of the
	// normal stdin/out/err. This is because due to the multi-process nature
	// of our plugin system, we can't use the normal process values so we
	// make our own custom one we pipe across.
	Stdout io.Reader
	Stderr io.Reader
}

// AppFunc creates app.App when they're requested from the server.
type AppFunc func() app.App

// Accept accepts connections on a listener and serves requests for
// each incoming connection. Accept blocks; the caller typically invokes
// it in a go statement.
func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("[ERR] plugin server: %s", err)
			return
		}

		go s.ServeConn(conn)
	}
}

// ServeConn runs a single connection.
//
// ServeConn blocks, serving the connection until the client hangs up.
func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	// First create the yamux server to wrap this connection
	mux, err := yamux.Server(conn, nil)
	if err != nil {
		conn.Close()
		log.Printf("[ERR] plugin: %s", err)
		return
	}

	// Accept the control connection
	control, err := mux.Accept()
	if err != nil {
		mux.Close()
		log.Printf("[ERR] plugin: %s", err)
		return
	}

	// Connect the stdstreams (in, out, err)
	stdstream := make([]net.Conn, 2)
	for i, _ := range stdstream {
		stdstream[i], err = mux.Accept()
		if err != nil {
			mux.Close()
			log.Printf("[ERR] plugin: accepting stream %d: %s", i, err)
			return
		}
	}

	// Copy std streams out to the proper place
	go copyStream("stdout", stdstream[0], s.Stdout)
	go copyStream("stderr", stdstream[1], s.Stderr)

	// Create the broker and start it up
	broker := newMuxBroker(mux)
	go broker.Run()

	// Use the control connection to build the dispenser and serve the
	// connection.
	server := rpc.NewServer()
	server.RegisterName("Dispenser", &dispenseServer{
		AppFunc: s.AppFunc,

		broker: broker,
	})
	server.ServeConn(control)
}

// dispenseServer dispenses variousinterface implementations for Terraform.
type dispenseServer struct {
	AppFunc AppFunc

	broker *muxBroker
}

func (d *dispenseServer) App(
	args interface{}, response *uint32) error {
	id := d.broker.NextId()
	*response = id

	go func() {
		conn, err := d.broker.Accept(id)
		if err != nil {
			log.Printf("[ERR] Plugin dispense: %s", err)
			return
		}

		serve(conn, "App", &AppServer{
			Broker: d.broker,
			App:    d.AppFunc(),
		})
	}()

	return nil
}

func acceptAndServe(mux *muxBroker, id uint32, n string, v interface{}) {
	conn, err := mux.Accept(id)
	if err != nil {
		log.Printf("[ERR] Plugin acceptAndServe: %s", err)
		return
	}

	serve(conn, n, v)
}

func serve(conn io.ReadWriteCloser, name string, v interface{}) {
	server := rpc.NewServer()
	if err := server.RegisterName(name, v); err != nil {
		log.Printf("[ERR] Plugin dispense: %s", err)
		return
	}

	server.ServeConn(conn)
}
