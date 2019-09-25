package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// IncDecNode represents a node with increment statement
type IncDecNode struct {
	*BaseNode
	LHS         Node
	IsIncrement bool
}

// Accept is part of visitor pattern.
func (node *IncDecNode) Accept(visitor Visitor) {
	visitor.VisitEnterIncDecNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveIncDecNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node.
func (node *IncDecNode) VisitChildren(visitor Visitor) {
	Accept(node.LHS, visitor)
}

func (node *IncDecNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		LHS      Node
		Typing   typing.Typing
	}{
		NodeType: "increment decrement statement",
		Token:    node.BaseNode.Tok,
		LHS:      node.LHS,
		Typing:   node.Typing,
	})
}

func CreateIncDecNode(tok *token.Token) *IncDecNode {
	var node IncDecNode
	node.BaseNode = CreateBaseNode(tok, nil)

	return &node
}
