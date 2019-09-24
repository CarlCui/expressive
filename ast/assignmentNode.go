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
	Identifier Node
	Expr       Node
	Operator   signature.Operator // if assignment is a compound assignment
}

// Accept is part of visitor pattern.
func (node *AssignmentNode) Accept(visitor Visitor) {
	visitor.VisitEnterAssignmentNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveAssignmentNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *AssignmentNode) VisitChildren(visitor Visitor) {
	Accept(node.Identifier, visitor)
	Accept(node.Expr, visitor)
}

func (node *AssignmentNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType   string
		Token      *token.Token
		Identifier Node
		Typing     typing.Typing
		Expr       Node
		Operator   signature.Operator
	}{
		NodeType:   "assignment",
		Token:      node.BaseNode.Tok,
		Identifier: node.Identifier,
		Typing:     node.Typing,
		Expr:       node.Expr,
		Operator:   node.Operator,
	})
}
