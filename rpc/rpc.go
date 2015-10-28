package rpc

import (
	"errors"
	"fmt"
	"net/rpc"
	"sync"

	"github.com/hashicorp/otto/app"
)

// nextId is the next ID to use for names registered.
var nextId uint32 = 0
var nextLock sync.Mutex

// Register registers an Otto thing with the RPC server and returns
// the name it is registered under.
func Register(server *rpc.Server, thing interface{}) (name string, err error) {
	nextLock.Lock()
	defer nextLock.Unlock()

	switch t := thing.(type) {
	case app.App:
		name = fmt.Sprintf("Otto%d", nextId)
		err = server.RegisterName(name, &AppServer{App: t})
	default:
		return "", errors.New("Unknown type to register for RPC server.")
	}

	nextId += 1
	return
}
