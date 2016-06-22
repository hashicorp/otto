package pongo2

import (
	"bytes"
)

type tagAutoescapeNode struct {
	wrapper    *NodeWrapper
	autoescape bool
}

func (node *tagAutoescapeNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	old := ctx.Autoescape
	ctx.Autoescape = node.autoescape

	err := node.wrapper.Execute(ctx, buffer)
	if err != nil {
		return err
	}

	ctx.Autoescape = old

	return nil
}

func tagAutoescapeParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	autoescape_node := &tagAutoescapeNode{}

	wrapper, _, err := doc.WrapUntilTag("endautoescape")
	if err != nil {
		return nil, err
	}
	autoescape_node.wrapper = wrapper

	mode_token := arguments.MatchType(TokenIdentifier)
	if mode_token == nil {
		return nil, arguments.Error("A mode is required for autoescape-tag.", nil)
	}
	if mode_token.Val == "on" {
		autoescape_node.autoescape = true
	} else if mode_token.Val == "off" {
		autoescape_node.autoescape = false
	} else {
		return nil, arguments.Error("Only 'on' or 'off' is valid as an autoescape-mode.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed autoescape-tag arguments.", nil)
	}

	return autoescape_node, nil
}

func init() {
	RegisterTag("autoescape", tagAutoescapeParser)
}
