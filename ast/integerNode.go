package ast

import "github.com/carlcui/expressive/token"

// IntegerNode represents an integer constant node.
type IntegerNode struct {
	BaseNode
	val int
}

// Accept is part of visitor pattern.
func (node *IntegerNode) Accept(visitor Visitor) {
	visitor.VisitIntegerNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *IntegerNode) VisitChildren(visitor Visitor) {

}

func (node *IntegerNode) Init(tok *token.Token) {
	node.BaseNode = BaseNode{tok: tok, parent: nil}

	// TODO: parse val
}
