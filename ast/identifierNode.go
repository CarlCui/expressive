package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
)

// IdentifierNode represents an identifier node.
type IdentifierNode struct {
	*BaseNode
}

// Accept is part of visitor pattern.
func (node *IdentifierNode) Accept(visitor Visitor) {
	visitor.VisitIdentifierNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *IdentifierNode) VisitChildren(visitor Visitor) {

}

func (node *IdentifierNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
	}{
		NodeType: "identifier",
		Token:    node.BaseNode.tok,
	})
}
