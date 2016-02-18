package plan

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

// EncodeHCL takes a list of plans and encodes it as HCL.
func EncodeHCL(ps []*Plan) *ast.File {
	items := make([]*ast.ObjectItem, 0, len(ps))
	for _, p := range ps {
		items = append(items, p.HCL())
	}

	return &ast.File{
		Node: &ast.ObjectList{Items: items},
	}
}

func (p *Plan) HCL() *ast.ObjectItem {
	items := make([]*ast.ObjectItem, 0, len(p.Tasks)+1)

	// Description if we have one
	if p.Description != "" {
		items = append(items, &ast.ObjectItem{
			Keys: []*ast.ObjectKey{
				&ast.ObjectKey{
					Token: token.Token{Type: token.IDENT, Text: "description"},
				},
			},
			Val: &ast.LiteralType{
				Token: token.Token{
					Type: token.STRING,
					Text: fmt.Sprintf(`"%s"`, p.Description),
				},
			},
			Assign: emptyAssign,
		})
	}

	// Inputs first if we have any
	if len(p.Inputs) > 0 {
		// Sort the keys so that it is deterministic and build the
		// list of items in the object.
		inputs := make([]*ast.ObjectItem, 0, len(p.Inputs))
		keys := make([]string, 0, len(p.Inputs))
		for k, _ := range p.Inputs {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			inputs = append(inputs, &ast.ObjectItem{
				Keys: []*ast.ObjectKey{
					&ast.ObjectKey{
						Token: token.Token{
							Type: token.STRING,
							Text: fmt.Sprintf(`"%s"`, k),
							Pos:  token.Pos{Line: len(inputs)},
						},
					},
				},
				Val: &ast.LiteralType{
					Token: token.Token{
						Type: token.STRING,
						Text: fmt.Sprintf(`"%s"`, p.Inputs[k]),
					},
				},
				Assign: emptyAssign,
			})
		}

		// Create the object itself and add it to our list of things
		items = append(items, &ast.ObjectItem{
			Keys: []*ast.ObjectKey{
				&ast.ObjectKey{
					Token: token.Token{Type: token.IDENT, Text: "inputs"},
				},
			},
			Val: &ast.ObjectType{
				List: &ast.ObjectList{
					Items: inputs,
				},
			},
		})
	}

	// Tasks
	for _, t := range p.Tasks {
		items = append(items, t.HCL())
	}

	// Return the plan
	return &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{Type: token.IDENT, Text: "plan"},
			},
		},
		Val: &ast.ObjectType{
			List: &ast.ObjectList{
				Items: items,
			},
		},
	}
}

func (t *Task) HCL() *ast.ObjectItem {
	items := make([]*ast.ObjectItem, 0, len(t.Args)+1)

	// Description if we have one
	if t.Description != "" {
		items = append(items, &ast.ObjectItem{
			Keys: []*ast.ObjectKey{
				&ast.ObjectKey{
					Token: token.Token{
						Type: token.IDENT,
						Text: "description",
						Pos:  token.Pos{Line: 1},
					},
				},
			},
			Val: &ast.LiteralType{
				Token: token.Token{
					Type: token.STRING,
					Text: fmt.Sprintf(`"%s"`, t.Description),
				},
			},
			Assign: emptyAssign,
		})
	}

	// Sort the args
	keys := make([]string, 0, len(t.Args))
	for k, _ := range t.Args {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// For each arg, add it to the object
	for _, k := range keys {
		arg := t.Args[k]

		var argToken token.Token
		switch v := arg.Value.(type) {
		case string:
			argToken.Type = token.STRING
			argToken.Text = fmt.Sprintf(`"%s"`, v)
		default:
			panic(fmt.Sprintf("Unknown arg type: %T", arg.Value))
		}

		items = append(items, &ast.ObjectItem{
			Keys: []*ast.ObjectKey{
				&ast.ObjectKey{
					Token: token.Token{
						Type: token.IDENT,
						Text: k,
						Pos:  token.Pos{Line: len(items)},
					},
				},
			},
			Val: &ast.LiteralType{
				Token: argToken,
			},
			Assign: emptyAssign,
		})
	}

	// Create the object
	return &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{Type: token.IDENT, Text: "task"},
			},
			&ast.ObjectKey{
				Token: token.Token{
					Type: token.STRING,
					Text: fmt.Sprintf(`"%s"`, t.Type),
				},
			},
		},
		Val: &ast.ObjectType{
			List: &ast.ObjectList{
				Items: items,
			},
		},
	}
}

var emptyAssign = token.Pos{Line: 1}
