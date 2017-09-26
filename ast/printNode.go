package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// PrintNode represents a node with print statement to stdout
type PrintNode struct {
	*BaseNode
	Expr Node
}

// Accept is part of visitor pattern.
func (node *PrintNode) Accept(visitor Visitor) {
	visitor.VisitEnterPrintNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeavePrintNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *PrintNode) VisitChildren(visitor Visitor) {
	Accept(node.Expr, visitor)
}

func (node *PrintNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Expr     Node
		Typing   typing.Typing
	}{
		NodeType: "print",
		Token:    node.BaseNode.Tok,
		Expr:     node.Expr,
		Typing:   node.Typing,
	})
}
