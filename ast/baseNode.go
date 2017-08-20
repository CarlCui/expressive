package ast

import (
	"github.com/carlcui/expressive/token"
)

// BaseNode has properties that all nodes have
type BaseNode struct {
	tok    *token.Token
	parent Node
}

// CreateBaseNode is a factory
func CreateBaseNode(tok *token.Token, parent Node) BaseNode {
	return BaseNode{tok: tok, parent: parent}
}

func (node *BaseNode) SetParent(parent Node) {
	node.parent = parent
}

func (node *BaseNode) GetParent() Node {
	return node.parent
}

// Accept is part of visitor pattern.
func (node *BaseNode) Accept(visitor Visitor) {
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *BaseNode) VisitChildren(visitor Visitor) {
}
