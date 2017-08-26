package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
)

// AssignmentNode represents a node with assignment statement
type AssignmentNode struct {
	*BaseNode
	Identifier Node
	Expr       Node
}

// Accept is part of visitor pattern.
func (node *AssignmentNode) Accept(visitor Visitor) {
	visitor.VisitEnterAssignmentNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveAssignmentNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *AssignmentNode) VisitChildren(visitor Visitor) {
	node.Identifier.Accept(visitor)
	node.Expr.Accept(visitor)
}

func (node *AssignmentNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType   string
		Token      *token.Token
		Identifier Node
		Expr       Node
	}{
		NodeType:   "assignment",
		Token:      node.BaseNode.tok,
		Identifier: node.Identifier,
		Expr:       node.Expr,
	})
}
