package pongo2

import (
	"bytes"
)

type tagIncludeNode struct {
	tpl                *Template
	filename_evaluator IEvaluator
	lazy               bool
	only               bool
	filename           string
	with_pairs         map[string]IEvaluator
}

func (node *tagIncludeNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	// Building the context for the template
	include_ctx := make(Context)

	// Fill the context with all data from the parent
	if !node.only {
		include_ctx.Update(ctx.Public)
		include_ctx.Update(ctx.Private)
	}

	// Put all custom with-pairs into the context
	for key, value := range node.with_pairs {
		val, err := value.Evaluate(ctx)
		if err != nil {
			return err
		}
		include_ctx[key] = val
	}

	// Execute the template
	if node.lazy {
		// Evaluate the filename
		filename, err := node.filename_evaluator.Evaluate(ctx)
		if err != nil {
			return err
		}

		if filename.String() == "" {
			return ctx.Error("Filename for 'include'-tag evaluated to an empty string.", nil)
		}

		// Get include-filename
		included_filename := ctx.template.set.resolveFilename(ctx.template, filename.String())

		included_tpl, err2 := ctx.template.set.FromFile(included_filename)
		if err2 != nil {
			return err2.(*Error)
		}
		err2 = included_tpl.ExecuteWriter(include_ctx, buffer)
		if err2 != nil {
			return err2.(*Error)
		}
		return nil
	} else {
		// Template is already parsed with static filename
		err := node.tpl.ExecuteWriter(include_ctx, buffer)
		if err != nil {
			return err.(*Error)
		}
		return nil
	}
}

func tagIncludeParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	include_node := &tagIncludeNode{
		with_pairs: make(map[string]IEvaluator),
	}

	if filename_token := arguments.MatchType(TokenString); filename_token != nil {
		// prepared, static template

		// Get include-filename
		included_filename := doc.template.set.resolveFilename(doc.template, filename_token.Val)

		// Parse the parent
		include_node.filename = included_filename
		included_tpl, err := doc.template.set.FromFile(included_filename)
		if err != nil {
			return nil, err.(*Error).updateFromTokenIfNeeded(doc.template, filename_token)
		}
		include_node.tpl = included_tpl
	} else {
		// No String, then the user wants to use lazy-evaluation (slower, but possible)
		filename_evaluator, err := arguments.ParseExpression()
		if err != nil {
			return nil, err.updateFromTokenIfNeeded(doc.template, filename_token)
		}
		include_node.filename_evaluator = filename_evaluator
		include_node.lazy = true
	}

	// After having parsed the filename we're gonna parse the with+only options
	if arguments.Match(TokenIdentifier, "with") != nil {
		for arguments.Remaining() > 0 {
			// We have at least one key=expr pair (because of starting "with")
			key_token := arguments.MatchType(TokenIdentifier)
			if key_token == nil {
				return nil, arguments.Error("Expected an identifier", nil)
			}
			if arguments.Match(TokenSymbol, "=") == nil {
				return nil, arguments.Error("Expected '='.", nil)
			}
			value_expr, err := arguments.ParseExpression()
			if err != nil {
				return nil, err.updateFromTokenIfNeeded(doc.template, key_token)
			}

			include_node.with_pairs[key_token.Val] = value_expr

			// Only?
			if arguments.Match(TokenIdentifier, "only") != nil {
				include_node.only = true
				break // stop parsing arguments because it's the last option
			}
		}
	}

	if arguments.Remaining() > 0 {
		return nil, arguments.Error("Malformed 'include'-tag arguments.", nil)
	}

	return include_node, nil
}

func init() {
	RegisterTag("include", tagIncludeParser)
}
