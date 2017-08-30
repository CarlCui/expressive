package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
)

// ErrorNode represents a node with syntax error.
type ErrorNode struct {
	*BaseNode
	Expected string
}

// Accept is part of visitor pattern.
func (node *ErrorNode) Accept(visitor Visitor) {
	visitor.VisitErrorNode(node)
}

// VisitChildren is part of visitor pattern. Error node does not have any child.
func (node *ErrorNode) VisitChildren(visitor Visitor) {

}

func (node *ErrorNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Message  string
	}{
		NodeType: "error",
		Token:    node.BaseNode.Tok,
		Message:  "expected " + node.Expected,
	})
}
