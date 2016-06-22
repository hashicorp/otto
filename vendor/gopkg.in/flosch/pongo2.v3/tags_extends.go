package pongo2

import (
	"bytes"
)

type tagExtendsNode struct {
	filename string
}

func (node *tagExtendsNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	return nil
}

func tagExtendsParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	extends_node := &tagExtendsNode{}

	if doc.template.level > 1 {
		return nil, arguments.Error("The 'extends' tag can only defined on root level.", start)
	}

	if doc.template.parent != nil {
		// Already one parent
		return nil, arguments.Error("This template has already one parent.", start)
	}

	if filename_token := arguments.MatchType(TokenString); filename_token != nil {
		// prepared, static template

		// Get parent's filename
		parent_filename := doc.template.set.resolveFilename(doc.template, filename_token.Val)

		// Parse the parent
		parent_template, err := doc.template.set.FromFile(parent_filename)
		if err != nil {
			return nil, err.(*Error)
		}

		// Keep track of things
		parent_template.child = doc.template
		doc.template.parent = parent_template
		extends_node.filename = parent_filename
	} else {
		return nil, arguments.Error("Tag 'extends' requires a template filename as string.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Tag 'extends' does only take 1 argument.", nil)
	}

	return extends_node, nil
}

func init() {
	RegisterTag("extends", tagExtendsParser)
}
