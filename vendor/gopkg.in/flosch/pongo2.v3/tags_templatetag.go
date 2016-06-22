package pongo2

import (
	"bytes"
)

type tagTemplateTagNode struct {
	content string
}

var templateTagMapping = map[string]string{
	"openblock":     "{%",
	"closeblock":    "%}",
	"openvariable":  "{{",
	"closevariable": "}}",
	"openbrace":     "{",
	"closebrace":    "}",
	"opencomment":   "{#",
	"closecomment":  "#}",
}

func (node *tagTemplateTagNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	buffer.WriteString(node.content)
	return nil
}

func tagTemplateTagParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	tt_node := &tagTemplateTagNode{}

	if arg_token := arguments.MatchType(TokenIdentifier); arg_token != nil {
		output, found := templateTagMapping[arg_token.Val]
		if !found {
			return nil, arguments.Error("Argument not found", arg_token)
		}
		tt_node.content = output
	} else {
		return nil, arguments.Error("Identifier expected.", nil)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed templatetag-tag argument.", nil)
	}

	return tt_node, nil
}

func init() {
	RegisterTag("templatetag", tagTemplateTagParser)
}
