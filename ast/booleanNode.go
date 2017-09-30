package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// BooleanNode represents an integer constant node.
type BooleanNode struct {
	*BaseNode
	val bool
}

// Accept is part of visitor pattern.
func (node *BooleanNode) Accept(visitor Visitor) {
	visitor.VisitBooleanNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *BooleanNode) VisitChildren(visitor Visitor) {

}

// Init initializes an integer node with a token
func (node *BooleanNode) Init(tok *token.Token) {
	node.BaseNode = CreateBaseNode(tok, nil)

	if tok.TokenType == token.TRUE {
		node.val = true
	} else {
		node.val = false
	}
}

func (node *BooleanNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Val      bool
		Typing   typing.Typing
	}{
		NodeType: "integer literal",
		Token:    node.BaseNode.Tok,
		Val:      node.val,
		Typing:   node.Typing,
	})
}
