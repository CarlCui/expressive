package ast

import (
	"encoding/json"
	"strconv"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// FloatNode represents a float constant node.
type FloatNode struct {
	*BaseNode
	Val float32
}

// Accept is part of visitor pattern.
func (node *FloatNode) Accept(visitor Visitor) {
	visitor.VisitFloatNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *FloatNode) VisitChildren(visitor Visitor) {

}

func (node *FloatNode) Init(tok *token.Token) {
	node.BaseNode = CreateBaseNode(tok, nil)

	val, err := strconv.ParseFloat(tok.Raw, 32)

	if err != nil {
		panic(tok.GetLocation() + ": error parsing float")
	}

	node.Val = float32(val)
}

func (node *FloatNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Val      float32
		Typing   typing.Typing
	}{
		NodeType: "float literal",
		Token:    node.BaseNode.Tok,
		Val:      node.Val,
		Typing:   node.Typing,
	})
}
