package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// TypeLiteralNode represents a node with a type literal
type TypeLiteralNode struct {
	*BaseNode
}

// Accept is part of visitor pattern.
func (node *TypeLiteralNode) Accept(visitor Visitor) {
	visitor.VisitTypeLiteralNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *TypeLiteralNode) VisitChildren(visitor Visitor) {
}

func (node *TypeLiteralNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Typing   typing.Typing
	}{
		NodeType: "type literal",
		Token:    node.BaseNode.Tok,
		Typing:   node.Typing,
	})
}
