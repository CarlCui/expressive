package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
)

// ArrayTypeNode represents a node with an array type
type ArrayTypeNode struct {
	*BaseNode
	SubType Node
}

// Accept is part of visitor pattern.
func (node *ArrayTypeNode) Accept(visitor Visitor) {
	visitor.VisitEnterArrayTypeNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveArrayTypeNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *ArrayTypeNode) VisitChildren(visitor Visitor) {
	Accept(node.SubType, visitor)
}

func (node *ArrayTypeNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		SubType  Node
	}{
		NodeType: "Array",
		SubType:  node.SubType,
	})
}

func CreateArrayTypeNode(tok *token.Token, subType Node) *ArrayTypeNode {
	var node ArrayTypeNode
	node.BaseNode = CreateBaseNode(tok, nil)

	node.SubType = subType

	return &node
}
