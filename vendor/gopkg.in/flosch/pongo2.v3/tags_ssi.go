package pongo2

import (
	"bytes"
	"io/ioutil"
)

type tagSSINode struct {
	filename string
	content  string
	template *Template
}

func (node *tagSSINode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	if node.template != nil {
		// Execute the template within the current context
		includeCtx := make(Context)
		includeCtx.Update(ctx.Public)
		includeCtx.Update(ctx.Private)

		err := node.template.ExecuteWriter(includeCtx, buffer)
		if err != nil {
			return err.(*Error)
		}
	} else {
		// Just print out the content
		buffer.WriteString(node.content)
	}
	return nil
}

func tagSSIParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	ssi_node := &tagSSINode{}

	if file_token := arguments.MatchType(TokenString); file_token != nil {
		ssi_node.filename = file_token.Val

		if arguments.Match(TokenIdentifier, "parsed") != nil {
			// parsed
			temporary_tpl, err := doc.template.set.FromFile(doc.template.set.resolveFilename(doc.template, file_token.Val))
			if err != nil {
				return nil, err.(*Error).updateFromTokenIfNeeded(doc.template, file_token)
			}
			ssi_node.template = temporary_tpl
		} else {
			// plaintext
			buf, err := ioutil.ReadFile(doc.template.set.resolveFilename(doc.template, file_token.Val))
			if err != nil {
				return nil, (&Error{
					Sender:   "tag:ssi",
					ErrorMsg: err.Error(),
				}).updateFromTokenIfNeeded(doc.template, file_token)
			}
			ssi_node.content = string(buf)
		}
	} else {
		return nil, arguments.Error("First argument must be a string.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed SSI-tag argument.", nil)
	}

	return ssi_node, nil
}

func init() {
	RegisterTag("ssi", tagSSIParser)
}
