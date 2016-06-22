package pongo2

import (
	"bytes"
)

// The root document
type nodeDocument struct {
	Nodes []INode
}

func (doc *nodeDocument) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	for _, n := range doc.Nodes {
		err := n.Execute(ctx, buffer)
		if err != nil {
			return err
		}
	}
	return nil
}
