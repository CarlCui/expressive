package ast

import (
	"github.com/carlcui/expressive/token"
)

// BaseNode has properties that all nodes have
type BaseNode struct {
	Tok    *token.Token
	Parent Node
}

// CreateBaseNode is a factory
func CreateBaseNode(tok *token.Token, parent Node) *BaseNode {
	return &BaseNode{Tok: tok, Parent: parent}
}

func (node *BaseNode) SetParent(parent Node) {
	node.Parent = parent
}

func (node *BaseNode) GetParent() Node {
	return node.Parent
}

// Accept is part of visitor pattern.
func (node *BaseNode) Accept(visitor Visitor) {
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *BaseNode) VisitChildren(visitor Visitor) {
}
