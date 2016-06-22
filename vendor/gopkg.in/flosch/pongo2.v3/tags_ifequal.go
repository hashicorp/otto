package pongo2

import (
	"bytes"
)

type tagIfEqualNode struct {
	var1, var2  IEvaluator
	thenWrapper *NodeWrapper
	elseWrapper *NodeWrapper
}

func (node *tagIfEqualNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	r1, err := node.var1.Evaluate(ctx)
	if err != nil {
		return err
	}
	r2, err := node.var2.Evaluate(ctx)
	if err != nil {
		return err
	}

	result := r1.EqualValueTo(r2)

	if result {
		return node.thenWrapper.Execute(ctx, buffer)
	} else {
		if node.elseWrapper != nil {
			return node.elseWrapper.Execute(ctx, buffer)
		}
	}
	return nil
}

func tagIfEqualParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	ifequal_node := &tagIfEqualNode{}

	// Parse two expressions
	var1, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	var2, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	ifequal_node.var1 = var1
	ifequal_node.var2 = var2

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("ifequal only takes 2 arguments.", nil)
	}

	// Wrap then/else-blocks
	wrapper, endargs, err := doc.WrapUntilTag("else", "endifequal")
	if err != nil {
		return nil, err
	}
	ifequal_node.thenWrapper = wrapper

	if endargs.Count() > 0 {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	if wrapper.Endtag == "else" {
		// if there's an else in the if-statement, we need the else-Block as well
		wrapper, endargs, err = doc.WrapUntilTag("endifequal")
		if err != nil {
			return nil, err
		}
		ifequal_node.elseWrapper = wrapper

		if endargs.Count() > 0 {
			return nil, endargs.Error("Arguments not allowed here.", nil)
		}
	}

	return ifequal_node, nil
}

func init() {
	RegisterTag("ifequal", tagIfEqualParser)
}
