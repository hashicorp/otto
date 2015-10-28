package rpc

import (
	"io"
	"log"
	"net/rpc"

	"github.com/hashicorp/otto/context"
)

// ContextSharedArgs is a struct that should be embedded directly into
// args structs that contain a context. It will be populated with the IDs
// that can be used to communicate back to the interfaces.
type ContextSharedArgs struct {
	UiId uint32
}

func connectContext(
	broker *muxBroker,
	ctx *context.Shared, args *ContextSharedArgs) (io.Closer, error) {
	closer := &multiCloser{}

	// Setup Ui
	conn, err := broker.Dial(args.UiId)
	if err != nil {
		return closer, err
	}
	client := rpc.NewClient(conn)
	closer.Closers = append(closer.Closers, client)
	ctx.Ui = &Ui{
		Client: client,
		Name:   "Ui",
	}

	return closer, nil
}

func serveContext(broker *muxBroker, ctx *context.Shared, args *ContextSharedArgs) {
	// Server the Ui
	id := broker.NextId()
	go acceptAndServe(broker, id, "Ui", &UiServer{
		Ui: ctx.Ui,
	})
	args.UiId = id

	// Set the context fields to nil so that they aren't sent over the
	// network (Go will just panic if we didn't do this).
	ctx.Directory = nil
	ctx.Ui = nil
}

// multiCloser is an io.Closer that closes multiple closers.
type multiCloser struct {
	Closers []io.Closer
}

func (c *multiCloser) Close() error {
	for _, single := range c.Closers {
		if err := single.Close(); err != nil {
			log.Printf("[ERR] rpc/context: close error: %s", err)
		}
	}

	return nil
}
