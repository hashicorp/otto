package pongo2

import (
	"bytes"
)

type tagIfNotEqualNode struct {
	var1, var2  IEvaluator
	thenWrapper *NodeWrapper
	elseWrapper *NodeWrapper
}

func (node *tagIfNotEqualNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	r1, err := node.var1.Evaluate(ctx)
	if err != nil {
		return err
	}
	r2, err := node.var2.Evaluate(ctx)
	if err != nil {
		return err
	}

	result := !r1.EqualValueTo(r2)

	if result {
		return node.thenWrapper.Execute(ctx, buffer)
	} else {
		if node.elseWrapper != nil {
			return node.elseWrapper.Execute(ctx, buffer)
		}
	}
	return nil
}

func tagIfNotEqualParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	ifnotequal_node := &tagIfNotEqualNode{}

	// Parse two expressions
	var1, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	var2, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	ifnotequal_node.var1 = var1
	ifnotequal_node.var2 = var2

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("ifequal only takes 2 arguments.", nil)
	}

	// Wrap then/else-blocks
	wrapper, endargs, err := doc.WrapUntilTag("else", "endifequal")
	if err != nil {
		return nil, err
	}
	ifnotequal_node.thenWrapper = wrapper

	if endargs.Count() > 0 {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	if wrapper.Endtag == "else" {
		// if there's an else in the if-statement, we need the else-Block as well
		wrapper, endargs, err = doc.WrapUntilTag("endifequal")
		if err != nil {
			return nil, err
		}
		ifnotequal_node.elseWrapper = wrapper

		if endargs.Count() > 0 {
			return nil, endargs.Error("Arguments not allowed here.", nil)
		}
	}

	return ifnotequal_node, nil
}

func init() {
	RegisterTag("ifnotequal", tagIfNotEqualParser)
}
