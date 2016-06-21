package pongo2

import (
	"bytes"
)

type tagFirstofNode struct {
	position *Token
	args     []IEvaluator
}

func (node *tagFirstofNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	for _, arg := range node.args {
		val, err := arg.Evaluate(ctx)
		if err != nil {
			return err
		}

		if val.IsTrue() {
			if ctx.Autoescape && !arg.FilterApplied("safe") {
				val, err = ApplyFilter("escape", val, nil)
				if err != nil {
					return err
				}
			}

			buffer.WriteString(val.String())
			return nil
		}
	}

	return nil
}

func tagFirstofParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	firstof_node := &tagFirstofNode{
		position: start,
	}

	for arguments.Remaining() > 0 {
		node, err := arguments.ParseExpression()
		if err != nil {
			return nil, err
		}
		firstof_node.args = append(firstof_node.args, node)
	}

	return firstof_node, nil
}

func init() {
	RegisterTag("firstof", tagFirstofParser)
}
