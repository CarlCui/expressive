package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
)

// BinaryOperatorNode represents a node with a binary operation (+, -, etc)
type BinaryOperatorNode struct {
	*BaseNode
	Lhs Node
	Rhs Node
}

// Accept is part of visitor pattern.
func (node *BinaryOperatorNode) Accept(visitor Visitor) {
	visitor.VisitEnterBinaryOepratorNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveBinaryOperatorNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *BinaryOperatorNode) VisitChildren(visitor Visitor) {
	node.Lhs.Accept(visitor)
	node.Rhs.Accept(visitor)
}

func (node *BinaryOperatorNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Lhs      Node
		Rhs      Node
	}{
		NodeType: "binary operator",
		Token:    node.BaseNode.tok,
		Lhs:      node.Lhs,
		Rhs:      node.Rhs,
	})
}

func CreateBinaryOperatorNode(tok *token.Token, lhs Node, rhs Node) Node {
	var node BinaryOperatorNode
	node.BaseNode = CreateBaseNode(tok, nil)

	lhs.SetParent(&node)
	rhs.SetParent(&node)

	node.Lhs = lhs
	node.Rhs = rhs

	return &node
}
