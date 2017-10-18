package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// PrintNode represents a node with print statement to stdout
type PrintNode struct {
	*BaseNode
	StringExpr Node
	Args       []Node
}

// Accept is part of visitor pattern.
func (node *PrintNode) Accept(visitor Visitor) {
	visitor.VisitEnterPrintNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeavePrintNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *PrintNode) VisitChildren(visitor Visitor) {
	Accept(node.StringExpr, visitor)

	for _, arg := range node.Args {
		Accept(arg, visitor)
	}
}

func (node *PrintNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType   string
		Token      *token.Token
		StringExpr Node
		Args       []Node
		Typing     typing.Typing
	}{
		NodeType:   "print",
		Token:      node.BaseNode.Tok,
		StringExpr: node.StringExpr,
		Args:       node.Args,
		Typing:     node.Typing,
	})
}
