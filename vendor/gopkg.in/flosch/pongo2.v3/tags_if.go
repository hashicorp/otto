package pongo2

import (
	"bytes"
)

type tagIfNode struct {
	conditions []IEvaluator
	wrappers   []*NodeWrapper
}

func (node *tagIfNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	for i, condition := range node.conditions {
		result, err := condition.Evaluate(ctx)
		if err != nil {
			return err
		}

		if result.IsTrue() {
			return node.wrappers[i].Execute(ctx, buffer)
		} else {
			// Last condition?
			if len(node.conditions) == i+1 && len(node.wrappers) > i+1 {
				return node.wrappers[i+1].Execute(ctx, buffer)
			}
		}
	}
	return nil
}

func tagIfParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	if_node := &tagIfNode{}

	// Parse first and main IF condition
	condition, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	if_node.conditions = append(if_node.conditions, condition)

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("If-condition is malformed.", nil)
	}

	// Check the rest
	for {
		wrapper, tag_args, err := doc.WrapUntilTag("elif", "else", "endif")
		if err != nil {
			return nil, err
		}
		if_node.wrappers = append(if_node.wrappers, wrapper)

		if wrapper.Endtag == "elif" {
			// elif can take a condition
			condition, err := tag_args.ParseExpression()
			if err != nil {
				return nil, err
			}
			if_node.conditions = append(if_node.conditions, condition)

			if tag_args.Remaining() > 0 {
				return nil, tag_args.Error("Elif-condition is malformed.", nil)
			}
		} else {
			if tag_args.Count() > 0 {
				// else/endif can't take any conditions
				return nil, tag_args.Error("Arguments not allowed here.", nil)
			}
		}

		if wrapper.Endtag == "endif" {
			break
		}
	}

	return if_node, nil
}

func init() {
	RegisterTag("if", tagIfParser)
}
