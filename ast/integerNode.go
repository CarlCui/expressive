package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
)

// IntegerNode represents an integer constant node.
type IntegerNode struct {
	*BaseNode
	val int
}

// Accept is part of visitor pattern.
func (node *IntegerNode) Accept(visitor Visitor) {
	visitor.VisitIntegerNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *IntegerNode) VisitChildren(visitor Visitor) {

}

// Init initializes an integer node with a token
func (node *IntegerNode) Init(tok *token.Token) {
	node.BaseNode = CreateBaseNode(tok, nil)

	// TODO: parse val
}

func (node *IntegerNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Val      int
	}{
		NodeType: "integer literal",
		Token:    node.BaseNode.Tok,
		Val:      node.val,
	})
}
