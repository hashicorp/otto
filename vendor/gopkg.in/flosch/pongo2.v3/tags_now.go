package pongo2

import (
	"bytes"
	"time"
)

type tagNowNode struct {
	position *Token
	format   string
	fake     bool
}

func (node *tagNowNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	var t time.Time
	if node.fake {
		t = time.Date(2014, time.February, 05, 18, 31, 45, 00, time.UTC)
	} else {
		t = time.Now()
	}

	buffer.WriteString(t.Format(node.format))

	return nil
}

func tagNowParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	now_node := &tagNowNode{
		position: start,
	}

	format_token := arguments.MatchType(TokenString)
	if format_token == nil {
		return nil, arguments.Error("Expected a format string.", nil)
	}
	now_node.format = format_token.Val

	if arguments.MatchOne(TokenIdentifier, "fake") != nil {
		now_node.fake = true
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed now-tag arguments.", nil)
	}

	return now_node, nil
}

func init() {
	RegisterTag("now", tagNowParser)
}
