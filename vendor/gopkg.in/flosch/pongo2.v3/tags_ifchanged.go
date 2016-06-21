package pongo2

import (
	"bytes"
)

type tagIfchangedNode struct {
	watched_expr []IEvaluator
	last_values  []*Value
	last_content []byte
	thenWrapper  *NodeWrapper
	elseWrapper  *NodeWrapper
}

func (node *tagIfchangedNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {

	if len(node.watched_expr) == 0 {
		// Check against own rendered body

		buf := bytes.NewBuffer(make([]byte, 0, 1024)) // 1 KiB
		err := node.thenWrapper.Execute(ctx, buf)
		if err != nil {
			return err
		}

		buf_bytes := buf.Bytes()
		if !bytes.Equal(node.last_content, buf_bytes) {
			// Rendered content changed, output it
			buffer.Write(buf_bytes)
			node.last_content = buf_bytes
		}
	} else {
		now_values := make([]*Value, 0, len(node.watched_expr))
		for _, expr := range node.watched_expr {
			val, err := expr.Evaluate(ctx)
			if err != nil {
				return err
			}
			now_values = append(now_values, val)
		}

		// Compare old to new values now
		changed := len(node.last_values) == 0

		for idx, old_val := range node.last_values {
			if !old_val.EqualValueTo(now_values[idx]) {
				changed = true
				break // we can stop here because ONE value changed
			}
		}

		node.last_values = now_values

		if changed {
			// Render thenWrapper
			err := node.thenWrapper.Execute(ctx, buffer)
			if err != nil {
				return err
			}
		} else {
			// Render elseWrapper
			err := node.elseWrapper.Execute(ctx, buffer)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func tagIfchangedParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	ifchanged_node := &tagIfchangedNode{}

	for arguments.Remaining() > 0 {
		// Parse condition
		expr, err := arguments.ParseExpression()
		if err != nil {
			return nil, err
		}
		ifchanged_node.watched_expr = append(ifchanged_node.watched_expr, expr)
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Ifchanged-arguments are malformed.", nil)
	}

	// Wrap then/else-blocks
	wrapper, endargs, err := doc.WrapUntilTag("else", "endifchanged")
	if err != nil {
		return nil, err
	}
	ifchanged_node.thenWrapper = wrapper

	if endargs.Count() > 0 {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	if wrapper.Endtag == "else" {
		// if there's an else in the if-statement, we need the else-Block as well
		wrapper, endargs, err = doc.WrapUntilTag("endifchanged")
		if err != nil {
			return nil, err
		}
		ifchanged_node.elseWrapper = wrapper

		if endargs.Count() > 0 {
			return nil, endargs.Error("Arguments not allowed here.", nil)
		}
	}

	return ifchanged_node, nil
}

func init() {
	RegisterTag("ifchanged", tagIfchangedParser)
}
