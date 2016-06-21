package pongo2

import (
	"bytes"
)

type tagForNode struct {
	key              string
	value            string // only for maps: for key, value in map
	object_evaluator IEvaluator
	reversed         bool

	bodyWrapper  *NodeWrapper
	emptyWrapper *NodeWrapper
}

type tagForLoopInformation struct {
	Counter     int
	Counter0    int
	Revcounter  int
	Revcounter0 int
	First       bool
	Last        bool
	Parentloop  *tagForLoopInformation
}

func (node *tagForNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) (forError *Error) {
	// Backup forloop (as parentloop in public context), key-name and value-name
	forCtx := NewChildExecutionContext(ctx)
	parentloop := forCtx.Private["forloop"]

	// Create loop struct
	loopInfo := &tagForLoopInformation{
		First: true,
	}

	// Is it a loop in a loop?
	if parentloop != nil {
		loopInfo.Parentloop = parentloop.(*tagForLoopInformation)
	}

	// Register loopInfo in public context
	forCtx.Private["forloop"] = loopInfo

	obj, err := node.object_evaluator.Evaluate(forCtx)
	if err != nil {
		return err
	}

	obj.IterateOrder(func(idx, count int, key, value *Value) bool {
		// There's something to iterate over (correct type and at least 1 item)

		// Update loop infos and public context
		forCtx.Private[node.key] = key
		if value != nil {
			forCtx.Private[node.value] = value
		}
		loopInfo.Counter = idx + 1
		loopInfo.Counter0 = idx
		if idx == 1 {
			loopInfo.First = false
		}
		if idx+1 == count {
			loopInfo.Last = true
		}
		loopInfo.Revcounter = count - idx        // TODO: Not sure about this, have to look it up
		loopInfo.Revcounter0 = count - (idx + 1) // TODO: Not sure about this, have to look it up

		// Render elements with updated context
		err := node.bodyWrapper.Execute(forCtx, buffer)
		if err != nil {
			forError = err
			return false
		}
		return true
	}, func() {
		// Nothing to iterate over (maybe wrong type or no items)
		if node.emptyWrapper != nil {
			err := node.emptyWrapper.Execute(forCtx, buffer)
			if err != nil {
				forError = err
			}
		}
	}, node.reversed)

	return nil
}

func tagForParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	for_node := &tagForNode{}

	// Arguments parsing
	var value_token *Token
	key_token := arguments.MatchType(TokenIdentifier)
	if key_token == nil {
		return nil, arguments.Error("Expected an key identifier as first argument for 'for'-tag", nil)
	}

	if arguments.Match(TokenSymbol, ",") != nil {
		// Value name is provided
		value_token = arguments.MatchType(TokenIdentifier)
		if value_token == nil {
			return nil, arguments.Error("Value name must be an identifier.", nil)
		}
	}

	if arguments.Match(TokenKeyword, "in") == nil {
		return nil, arguments.Error("Expected keyword 'in'.", nil)
	}

	object_evaluator, err := arguments.ParseExpression()
	if err != nil {
		return nil, err
	}
	for_node.object_evaluator = object_evaluator
	for_node.key = key_token.Val
	if value_token != nil {
		for_node.value = value_token.Val
	}

	if arguments.MatchOne(TokenIdentifier, "reversed") != nil {
		for_node.reversed = true
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed for-loop arguments.", nil)
	}

	// Body wrapping
	wrapper, endargs, err := doc.WrapUntilTag("empty", "endfor")
	if err != nil {
		return nil, err
	}
	for_node.bodyWrapper = wrapper

	if endargs.Count() > 0 {
		return nil, endargs.Error("Arguments not allowed here.", nil)
	}

	if wrapper.Endtag == "empty" {
		// if there's an else in the if-statement, we need the else-Block as well
		wrapper, endargs, err = doc.WrapUntilTag("endfor")
		if err != nil {
			return nil, err
		}
		for_node.emptyWrapper = wrapper

		if endargs.Count() > 0 {
			return nil, endargs.Error("Arguments not allowed here.", nil)
		}
	}

	return for_node, nil
}

func init() {
	RegisterTag("for", tagForParser)
}
