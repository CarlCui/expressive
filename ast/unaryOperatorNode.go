package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// UnaryOperatorNode represents a node with a left-associated unary operator
type UnaryOperatorNode struct {
	*BaseNode
	Expr     Node
	Operator signature.Operator
}

// Accept is part of visitor pattern.
func (node *UnaryOperatorNode) Accept(visitor Visitor) {
	visitor.VisitEnterUnaryOperatorNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveUnaryOperatorNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *UnaryOperatorNode) VisitChildren(visitor Visitor) {
	Accept(node.Expr, visitor)
}

func CreateUnaryOperatorNode(tok *token.Token, operator signature.Operator, expr Node) Node {
	var node UnaryOperatorNode
	node.BaseNode = CreateBaseNode(tok, nil)

	expr.SetParent(&node)

	node.Expr = expr

	node.Operator = operator

	return &node
}

func (node *UnaryOperatorNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Typing   typing.Typing
		Rhs      Node
	}{
		NodeType: "unary operator",
		Token:    node.BaseNode.Tok,
		Typing:   node.Typing,
		Rhs:      node.Expr,
	})
}
