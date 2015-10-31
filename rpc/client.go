package rpc

import (
	"io"
	"net"
	"net/rpc"

	"github.com/hashicorp/otto/app"
	"github.com/hashicorp/yamux"
)

// Client connects to a Server in order to request plugin implementations
// for Otto.
type Client struct {
	broker  *muxBroker
	control *rpc.Client
}

// Dial opens a connection to an RPC server and returns a client.
func Dial(network, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		// Make sure to set keep alive so that the connection doesn't die
		tcpConn.SetKeepAlive(true)
	}

	return NewClient(conn)
}

// NewClient creates a client from an already-open connection-like value.
// Dial is typically used instead.
func NewClient(conn io.ReadWriteCloser) (*Client, error) {
	// Create the yamux client so we can multiplex
	mux, err := yamux.Client(conn, nil)
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Connect to the control stream.
	control, err := mux.Open()
	if err != nil {
		mux.Close()
		return nil, err
	}

	// Create the broker and start it up
	broker := newMuxBroker(mux)
	go broker.Run()

	// Build the client using our broker and control channel.
	return &Client{
		broker:  broker,
		control: rpc.NewClient(control),
	}, nil
}

// Close closes the connection. The client is no longer usable after this
// is called.
func (c *Client) Close() error {
	if err := c.control.Close(); err != nil {
		return err
	}

	return c.broker.Close()
}

func (c *Client) App() (app.App, error) {
	var id uint32
	if err := c.control.Call(
		"Dispenser.App", new(interface{}), &id); err != nil {
		return nil, err
	}

	conn, err := c.broker.Dial(id)
	if err != nil {
		return nil, err
	}

	return &App{
		Broker: c.broker,
		Client: rpc.NewClient(conn),
		Name:   "App",
	}, nil
}
