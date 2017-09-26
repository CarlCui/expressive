package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// BinaryOperatorNode represents a node with a binary operation (+, -, etc)
type BinaryOperatorNode struct {
	*BaseNode
	Lhs      Node
	Rhs      Node
	Operator signature.Operator
}

// Accept is part of visitor pattern.
func (node *BinaryOperatorNode) Accept(visitor Visitor) {
	visitor.VisitEnterBinaryOepratorNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveBinaryOperatorNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *BinaryOperatorNode) VisitChildren(visitor Visitor) {
	Accept(node.Lhs, visitor)
	Accept(node.Rhs, visitor)
}

func (node *BinaryOperatorNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Typing   typing.Typing
		Lhs      Node
		Rhs      Node
	}{
		NodeType: "binary operator",
		Token:    node.BaseNode.Tok,
		Typing:   node.Typing,
		Lhs:      node.Lhs,
		Rhs:      node.Rhs,
	})
}

func CreateBinaryOperatorNode(tok *token.Token, operator signature.Operator, lhs Node, rhs Node) Node {
	var node BinaryOperatorNode
	node.BaseNode = CreateBaseNode(tok, nil)

	lhs.SetParent(&node)
	rhs.SetParent(&node)

	node.Lhs = lhs
	node.Rhs = rhs

	node.Operator = operator

	return &node
}
