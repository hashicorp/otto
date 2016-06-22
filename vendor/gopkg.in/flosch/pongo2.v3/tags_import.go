package pongo2

import (
	"bytes"
	"fmt"
)

type tagImportNode struct {
	position *Token
	filename string
	template *Template
	macros   map[string]*tagMacroNode // alias/name -> macro instance
}

func (node *tagImportNode) Execute(ctx *ExecutionContext, buffer *bytes.Buffer) *Error {
	for name, macro := range node.macros {
		func(name string, macro *tagMacroNode) {
			ctx.Private[name] = func(args ...*Value) *Value {
				return macro.call(ctx, args...)
			}
		}(name, macro)
	}
	return nil
}

func tagImportParser(doc *Parser, start *Token, arguments *Parser) (INodeTag, *Error) {
	import_node := &tagImportNode{
		position: start,
		macros:   make(map[string]*tagMacroNode),
	}

	filename_token := arguments.MatchType(TokenString)
	if filename_token == nil {
		return nil, arguments.Error("Import-tag needs a filename as string.", nil)
	}

	import_node.filename = doc.template.set.resolveFilename(doc.template, filename_token.Val)

	if arguments.Remaining() == 0 {
		return nil, arguments.Error("You must at least specify one macro to import.", nil)
	}

	// Compile the given template
	tpl, err := doc.template.set.FromFile(import_node.filename)
	if err != nil {
		return nil, err.(*Error).updateFromTokenIfNeeded(doc.template, start)
	}

	for arguments.Remaining() > 0 {
		macro_name_token := arguments.MatchType(TokenIdentifier)
		if macro_name_token == nil {
			return nil, arguments.Error("Expected macro name (identifier).", nil)
		}

		as_name := macro_name_token.Val
		if arguments.Match(TokenKeyword, "as") != nil {
			alias_token := arguments.MatchType(TokenIdentifier)
			if alias_token == nil {
				return nil, arguments.Error("Expected macro alias name (identifier).", nil)
			}
			as_name = alias_token.Val
		}

		macro_instance, has := tpl.exported_macros[macro_name_token.Val]
		if !has {
			return nil, arguments.Error(fmt.Sprintf("Macro '%s' not found (or not exported) in '%s'.", macro_name_token.Val,
				import_node.filename), macro_name_token)
		}

		import_node.macros[as_name] = macro_instance

		if arguments.Remaining() == 0 {
			break
		}

		if arguments.Match(TokenSymbol, ",") == nil {
			return nil, arguments.Error("Expected ','.", nil)
		}
	}

	return import_node, nil
}

func init() {
	RegisterTag("import", tagImportParser)
}
