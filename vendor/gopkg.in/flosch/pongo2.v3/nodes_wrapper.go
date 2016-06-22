package pongo2

import (
	"bytes"
)

type NodeWrapper struct {
	Endtag string
	nodes  []INode
}

func (wrapper *NodeWrapper) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	for _, n := range wrapper.nodes {
		err := n.Execute(ctx, buffer)
		if err != nil {
			return err
		}
	}
	return nil
}
