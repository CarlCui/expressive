package ast

import (
	"encoding/json"
	"strconv"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// IntegerNode represents an integer constant node.
type IntegerNode struct {
	*BaseNode
	Val int
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

	val, err := strconv.Atoi(tok.Raw)

	if err != nil {
		panic(tok.GetLocation() + ": error parsing int")
	}

	node.Val = val
}

func (node *IntegerNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Val      int
		Typing   typing.Typing
	}{
		NodeType: "integer literal",
		Token:    node.BaseNode.Tok,
		Val:      node.Val,
		Typing:   node.Typing,
	})
}
