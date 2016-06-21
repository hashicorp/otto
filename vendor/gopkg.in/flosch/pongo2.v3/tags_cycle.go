package pongo2

import (
	"bytes"
)

type tagCycleValue struct {
	node  *tagCycleNode
	value *Value
}

type tagCycleNode struct {
	position *Token
	args     []IEvaluator
	idx      int
	as_name  string
	silent   bool
}

func (cv *tagCycleValue) String() string {
	return cv.value.String()
}

func (node *tagCycleNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	item := node.args[node.idx%len(node.args)]
	node.idx++

	val, err := item.Evaluate(ctx)
	if err != nil {
		return err
	}

	if t, ok := val.Interface().(*tagCycleValue); ok {
		// {% cycle "test1" "test2"
		// {% cycle cycleitem %}

		// Update the cycle value with next value
		item := t.node.args[t.node.idx%len(t.node.args)]
		t.node.idx++

		val, err := item.Evaluate(ctx)
		if err != nil {
			return err
		}

		t.value = val

		if !t.node.silent {
			buffer.WriteString(val.String())
		}
	} else {
		// Regular call

		cycle_value := &tagCycleValue{
			node:  node,
			value: val,
		}

		if node.as_name != "" {
			ctx.Private[node.as_name] = cycle_value
		}
		if !node.silent {
			buffer.WriteString(val.String())
		}
	}

	return nil
}

// HINT: We're not supporting the old comma-seperated list of expresions argument-style
func tagCycleParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	cycle_node := &tagCycleNode{
		position: start,
	}

	for arguments.Remaining() > 0 {
		node, err := arguments.ParseExpression()
		if err != nil {
			return nil, err
		}
		cycle_node.args = append(cycle_node.args, node)

		if arguments.MatchOne(TokenKeyword, "as") != nil {
			// as

			name_token := arguments.MatchType(TokenIdentifier)
			if name_token == nil {
				return nil, arguments.Error("Name (identifier) expected after 'as'.", nil)
			}
			cycle_node.as_name = name_token.Val

			if arguments.MatchOne(TokenIdentifier, "silent") != nil {
				cycle_node.silent = true
			}

			// Now we're finished
			break
		}
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed cycle-tag.", nil)
	}

	return cycle_node, nil
}

func init() {
	RegisterTag("cycle", tagCycleParser)
}
