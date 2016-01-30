// Package ast declares the types used to represent syntax trees for plans
package ast

import (
	"fmt"

	"github.com/hashicorp/otto/plan"
	"github.com/hashicorp/otto/plan/lang/token"
)

// Node is an element in the abstract syntax tree.
type Node interface {
	node()
	Pos() token.Pos
}

type ObjectList struct{}
type ObjectItem struct{}

func (File) node()     {}
func (PlanList) node() {}
func (Plan) node()     {}

func (ObjectKey) node()    {}
func (Comment) node()      {}
func (CommentGroup) node() {}
func (ObjectType) node()   {}
func (LiteralType) node()  {}
func (ListType) node()     {}

// File represents a single HCL file
type File struct {
	Node     Node            // usually a *PlanList
	Comments []*CommentGroup // list of all comments in the source
}

func (f *File) Plans() []*plan.Plan {
	pl, ok := f.Node.(*PlanList)
	if !ok {
		return nil
	}

	return nil
}

func (f *File) Pos() token.Pos {
	return f.Node.Pos()
}

// PlanList represents a list of plans.
type PlanList struct {
	Items []*Plan
}

func (n *PlanList) Add(item *Plan) {
	n.Items = append(n.Items, item)
}

func (n *PlanList) Pos() token.Pos {
	// always returns the uninitiliazed position
	return n.Items[0].Pos()
}

// Plan represents a single plan.
type Plan struct {
	Token token.Token // token for "Plan"

	Lbrace token.Pos // position of "{"
	Rbrace token.Pos // position of "}"

	// TODO: metadata

	LeadComment *CommentGroup // associated lead comment
	LineComment *CommentGroup // associated line comment
}

func (n *Plan) Pos() token.Pos {
	return n.Token.Pos
}

// ObjectKeys are either an identifier or of type string.
type ObjectKey struct {
	Token token.Token
}

func (o *ObjectKey) Pos() token.Pos {
	return o.Token.Pos
}

// LiteralType represents a literal of basic type. Valid types are:
// token.NUMBER, token.FLOAT, token.BOOL and token.STRING
type LiteralType struct {
	Token token.Token

	// associated line comment, only when used in a list
	LineComment *CommentGroup
}

func (l *LiteralType) Pos() token.Pos {
	return l.Token.Pos
}

// ListStatement represents a HCL List type
type ListType struct {
	Lbrack token.Pos // position of "["
	Rbrack token.Pos // position of "]"
	List   []Node    // the elements in lexical order
}

func (l *ListType) Pos() token.Pos {
	return l.Lbrack
}

func (l *ListType) Add(node Node) {
	l.List = append(l.List, node)
}

// ObjectType represents a HCL Object Type
type ObjectType struct {
	Lbrace token.Pos   // position of "{"
	Rbrace token.Pos   // position of "}"
	List   *ObjectList // the nodes in lexical order
}

func (o *ObjectType) Pos() token.Pos {
	return o.Lbrace
}

// Comment node represents a single //, # style or /*- style commment
type Comment struct {
	Start token.Pos // position of / or #
	Text  string
}

func (c *Comment) Pos() token.Pos {
	return c.Start
}

// CommentGroup node represents a sequence of comments with no other tokens and
// no empty lines between.
type CommentGroup struct {
	List []*Comment // len(List) > 0
}

func (c *CommentGroup) Pos() token.Pos {
	return c.List[0].Pos()
}

//-------------------------------------------------------------------
// GoStringer
//-------------------------------------------------------------------

func (o *ObjectKey) GoString() string { return fmt.Sprintf("*%#v", *o) }
