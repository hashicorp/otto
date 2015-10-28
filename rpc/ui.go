package rpc

import (
	"log"
	"net/rpc"

	"github.com/hashicorp/otto/ui"
)

// Ui is an implementatin of ui.Ui that communicates over RPC.
type Ui struct {
	Client *rpc.Client
	Name   string
}

func (i *Ui) Header(msg string)  { i.basicCall("Header", msg) }
func (i *Ui) Message(msg string) { i.basicCall("Message", msg) }
func (i *Ui) Raw(msg string)     { i.basicCall("Raw", msg) }

func (i *Ui) Input(opts *ui.InputOpts) (string, error) {
	var resp UiInputResponse
	err := i.Client.Call(i.Name+".Input", opts, &resp)
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		err = resp.Error
		return "", err
	}

	return resp.Value, nil
}

func (i *Ui) basicCall(kind, msg string) {
	err := i.Client.Call(i.Name+"."+kind, msg, nil)
	if err != nil {
		log.Printf("[ERR] rpc/ui: %s", err)
	}
}

type UiInputResponse struct {
	Value string
	Error *BasicError
}

// UiServer is a net/rpc compatible structure for serving
// a Ui. This should not be used directly.
type UiServer struct {
	Ui ui.Ui
}

func (s *UiServer) Header(msg string, reply *interface{}) error {
	s.Ui.Header(msg)
	return nil
}

func (s *UiServer) Message(msg string, reply *interface{}) error {
	s.Ui.Message(msg)
	return nil
}

func (s *UiServer) Raw(msg string, reply *interface{}) error {
	s.Ui.Raw(msg)
	return nil
}

func (s *UiServer) Input(
	opts *ui.InputOpts,
	reply *UiInputResponse) error {
	value, err := s.Ui.Input(opts)
	*reply = UiInputResponse{
		Value: value,
		Error: NewBasicError(err),
	}

	return nil
}
