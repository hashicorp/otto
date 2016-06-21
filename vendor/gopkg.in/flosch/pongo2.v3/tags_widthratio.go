package pongo2

import (
	"bytes"
	"fmt"
	"math"
)

type tagWidthratioNode struct {
	position     *Token
	current, max IEvaluator
	width        IEvaluator
	ctx_name     string
}

func (node *tagWidthratioNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	current, err := node.current.Evaluate(ctx)
	if err != nil {
		return err
	}

	max, err := node.max.Evaluate(ctx)
	if err != nil {
		return err
	}

	width, err := node.width.Evaluate(ctx)
	if err != nil {
		return err
	}

	value := int(math.Ceil(current.Float()/max.Float()*width.Float() + 0.5))

	if node.ctx_name == "" {
		buffer.WriteString(fmt.Sprintf("%d", value))
	} else {
		ctx.Private[node.ctx_name] = value
	}

	return nil
}

func tagWidthratioParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	widthratio_node := &tagWidthratioNode{
		position: start,
	}

	current, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	widthratio_node.current = current

	max, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	widthratio_node.max = max

	width, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	widthratio_node.width = width

	if arguments.MatchOne(TokenKeyword, "as") != nil {
		// Name follows
		name_token := arguments.MatchType(TokenIdentifier)
		if name_token == nil {
			return nil, arguments.Error("Expected name (identifier).", nil)
		}
		widthratio_node.ctx_name = name_token.Val
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed widthratio-tag arguments.", nil)
	}

	return widthratio_node, nil
}

func init() {
	RegisterTag("widthratio", tagWidthratioParser)
}
