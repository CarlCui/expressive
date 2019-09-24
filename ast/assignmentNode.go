package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// AssignmentNode represents a node with assignment statement
type AssignmentNode struct {
	*BaseNode
	LHS      Node
	RHS      Node
	Operator signature.Operator // if assignment is a compound assignment
}

// Accept is part of visitor pattern.
func (node *AssignmentNode) Accept(visitor Visitor) {
	visitor.VisitEnterAssignmentNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveAssignmentNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *AssignmentNode) VisitChildren(visitor Visitor) {
	Accept(node.LHS, visitor)
	Accept(node.RHS, visitor)
}

func (node *AssignmentNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		LHS      Node
		Typing   typing.Typing
		RHS      Node
		Operator signature.Operator
	}{
		NodeType: "assignment",
		Token:    node.BaseNode.Tok,
		LHS:      node.LHS,
		Typing:   node.Typing,
		RHS:      node.RHS,
		Operator: node.Operator,
	})
}

func CreateAssignmentStmtNode(tok *token.Token) *AssignmentNode {
	var node AssignmentNode
	node.BaseNode = CreateBaseNode(tok, nil)

	return &node
}
