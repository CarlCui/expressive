package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// TernaryOperatorNode represents a node with a binary operation (+, -, etc)
type TernaryOperatorNode struct {
	*BaseNode
	Expr1    Node
	Expr2    Node
	Expr3    Node
	Operator signature.Operator
}

// Accept is part of visitor pattern.
func (node *TernaryOperatorNode) Accept(visitor Visitor) {
	visitor.VisitEnterTernaryOperatorNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveTernaryOperatorNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *TernaryOperatorNode) VisitChildren(visitor Visitor) {
	Accept(node.Expr1, visitor)
	Accept(node.Expr2, visitor)
	Accept(node.Expr3, visitor)
}

func CreateTernaryOperatorNode(tok *token.Token, operator signature.Operator, expr1 Node, expr2 Node, expr3 Node) Node {
	var node TernaryOperatorNode
	node.BaseNode = CreateBaseNode(tok, nil)

	expr1.SetParent(&node)
	expr2.SetParent(&node)
	expr3.SetParent(&node)

	node.Expr1 = expr1
	node.Expr2 = expr2
	node.Expr3 = expr3

	node.Operator = operator

	return &node
}

func (node *TernaryOperatorNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Typing   typing.Typing
		Expr1    Node
		Expr2    Node
		Expr3    Node
	}{
		NodeType: "ternary operator",
		Token:    node.BaseNode.Tok,
		Typing:   node.Typing,
		Expr1:    node.Expr1,
		Expr2:    node.Expr2,
		Expr3:    node.Expr3,
	})
}
