package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
)

// VariableDeclarationNode represents a node with variable declaration statement
type VariableDeclarationNode struct {
	*BaseNode
	Identifier   Node
	DeclaredType Node
	Expr         Node
}

// Accept is part of visitor pattern.
func (node *VariableDeclarationNode) Accept(visitor Visitor) {
	visitor.VisitEnterVariableDeclarationNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveVariableDeclarationNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *VariableDeclarationNode) VisitChildren(visitor Visitor) {
	Accept(node.Identifier, visitor)
	Accept(node.DeclaredType, visitor)
	Accept(node.Expr, visitor)
}

func (node *VariableDeclarationNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType     string
		Token        *token.Token
		Identifier   Node
		DeclaredType Node
		Expr         Node
	}{
		NodeType:     "variable declaration",
		Token:        node.BaseNode.Tok,
		Identifier:   node.Identifier,
		DeclaredType: node.DeclaredType,
		Expr:         node.Expr,
	})
}
