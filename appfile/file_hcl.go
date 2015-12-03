package appfile

import (
	"fmt"

	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/hashicorp/hcl/hcl/token"
)

var emptyAssign = token.Pos{Line: 1}

// HCL converts the Appfile to an HCL AST, allowing printing of the Appfile
// back to HCL.
//
// Note that if you parsed the File from HCL, this will not convert it
// back to the same HCL. Comments, in particular, won't be preserved.
func (f *File) HCL() *ast.File {
	// Convert all the various components into members of the root object
	items := make([]*ast.ObjectItem, 0, 10+len(f.Imports))
	for _, imp := range f.Imports {
		items = append(items, imp.HCL())
	}
	items = append(items, f.Application.HCL())
	items = append(items, f.Project.HCL())
	for _, infra := range f.Infrastructure {
		items = append(items, infra.HCL())
	}
	items = append(items, f.Customization.HCL()...)

	// Finalize
	return &ast.File{
		Node: &ast.ObjectList{
			Items: items,
		},
	}
}

func (f *CustomizationSet) HCL() []*ast.ObjectItem {
	if f == nil {
		return nil
	}

	items := make([]*ast.ObjectItem, 0, len(f.Raw))
	for _, c := range f.Raw {
		items = append(items, c.HCL())
	}

	return items
}

func (f *Customization) HCL() *ast.ObjectItem {
	items := make([]*ast.ObjectItem, 0, len(f.Config))
	for k, v := range f.Config {
		var val ast.Node
		switch t := v.(type) {
		case string:
			val = &ast.LiteralType{
				Token: token.Token{
					Type: token.STRING,
					Text: fmt.Sprintf(`"%s"`, t),
				},
			}
		default:
			panic(fmt.Sprintf("can't convert to HCL: %T", t))
		}

		items = append(items, &ast.ObjectItem{
			Keys: []*ast.ObjectKey{
				&ast.ObjectKey{
					Token: token.Token{Type: token.IDENT, Text: k},
				},
			},
			Val:    val,
			Assign: emptyAssign,
		})
	}

	return &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{Type: token.IDENT, Text: "customization"},
			},
			&ast.ObjectKey{
				Token: token.Token{
					Type: token.STRING,
					Text: fmt.Sprintf(`"%s"`, f.Type),
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

func (f *Dependency) HCL() *ast.ObjectItem {
	items := make([]*ast.ObjectItem, 0, 1)
	items = append(items, &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{Type: token.IDENT, Text: "source"},
			},
		},
		Val: &ast.LiteralType{
			Token: token.Token{
				Type: token.STRING,
				Text: fmt.Sprintf(`"%s"`, f.Source),
			},
		},
		Assign: emptyAssign,
	})

	return &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{Type: token.IDENT, Text: "dependency"},
			},
		},
		Val: &ast.ObjectType{
			List: &ast.ObjectList{
				Items: items,
			},
		},
	}
}

func (f *Import) HCL() *ast.ObjectItem {
	return &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{Type: token.IDENT, Text: "import"},
			},
			&ast.ObjectKey{
				Token: token.Token{
					Type: token.STRING,
					Text: fmt.Sprintf(`"%s"`, f.Source),
				},
			},
		},
		Val: &ast.ObjectType{},
	}
}

func (f *Application) HCL() *ast.ObjectItem {
	items := make([]*ast.ObjectItem, 0, 2+len(f.Dependencies))
	items = append(items, &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{
					Type: token.IDENT,
					Text: "name",
					Pos:  token.Pos{Line: 1},
				},
			},
		},
		Val: &ast.LiteralType{
			Token: token.Token{
				Type: token.STRING,
				Text: fmt.Sprintf(`"%s"`, f.Name),
			},
		},
		Assign: emptyAssign,
	})
	if f.Type != "" {
		items = append(items, &ast.ObjectItem{
			Keys: []*ast.ObjectKey{
				&ast.ObjectKey{
					Token: token.Token{
						Type: token.IDENT,
						Text: "type",
						Pos:  token.Pos{Line: 2},
					},
				},
			},
			Val: &ast.LiteralType{
				Token: token.Token{
					Type: token.STRING,
					Text: fmt.Sprintf(`"%s"`, f.Type),
				},
			},
			Assign: emptyAssign,
		})
	}
	for _, dep := range f.Dependencies {
		item := dep.HCL()
		items = append(items, item)
	}

	return &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{Type: token.IDENT, Text: "application"},
			},
		},
		Val: &ast.ObjectType{
			List: &ast.ObjectList{
				Items: items,
			},
		},
	}
}

func (f *Project) HCL() *ast.ObjectItem {
	items := make([]*ast.ObjectItem, 0, 2)
	items = append(items, &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{
					Type: token.IDENT,
					Text: "name",
					Pos:  token.Pos{Line: 1},
				},
			},
		},
		Val: &ast.LiteralType{
			Token: token.Token{
				Type: token.STRING,
				Text: fmt.Sprintf(`"%s"`, f.Name),
			},
		},
		Assign: emptyAssign,
	})
	items = append(items, &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{
					Type: token.IDENT,
					Text: "infrastructure",
					Pos:  token.Pos{Line: 2},
				},
			},
		},
		Val: &ast.LiteralType{
			Token: token.Token{
				Type: token.STRING,
				Text: fmt.Sprintf(`"%s"`, f.Infrastructure),
			},
		},
		Assign: emptyAssign,
	})

	return &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{Type: token.IDENT, Text: "project"},
			},
		},
		Val: &ast.ObjectType{
			List: &ast.ObjectList{
				Items: items,
			},
		},
	}
}

func (f *Infrastructure) HCL() *ast.ObjectItem {
	items := make([]*ast.ObjectItem, 0, 3+len(f.Foundations))
	items = append(items, &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{
					Type: token.IDENT,
					Text: "name",
					Pos:  token.Pos{Line: 1},
				},
			},
		},
		Val: &ast.LiteralType{
			Token: token.Token{
				Type: token.STRING,
				Text: fmt.Sprintf(`"%s"`, f.Name),
			},
		},
		Assign: emptyAssign,
	})
	items = append(items, &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{
					Type: token.IDENT,
					Text: "type",
					Pos:  token.Pos{Line: 2},
				},
			},
		},
		Val: &ast.LiteralType{
			Token: token.Token{
				Type: token.STRING,
				Text: fmt.Sprintf(`"%s"`, f.Type),
			},
		},
		Assign: emptyAssign,
	})
	items = append(items, &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{
					Type: token.IDENT,
					Text: "flavor",
					Pos:  token.Pos{Line: 3},
				},
			},
		},
		Val: &ast.LiteralType{
			Token: token.Token{
				Type: token.STRING,
				Text: fmt.Sprintf(`"%s"`, f.Flavor),
			},
		},
		Assign: emptyAssign,
	})

	return &ast.ObjectItem{
		Keys: []*ast.ObjectKey{
			&ast.ObjectKey{
				Token: token.Token{Type: token.IDENT, Text: "infrastructure"},
			},
		},
		Val: &ast.ObjectType{
			List: &ast.ObjectList{
				Items: items,
			},
		},
	}
}
