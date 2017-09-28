package ast

import (
	"github.com/carlcui/expressive/symbolTable"
	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// BaseNode has properties that all nodes have
type BaseNode struct {
	Tok    *token.Token
	Parent Node
	Typing typing.Typing
	Scope  *symbolTable.Scope
}

// CreateBaseNode is a factory
func CreateBaseNode(tok *token.Token, parent Node) *BaseNode {
	return &BaseNode{Tok: tok, Parent: parent}
}

func (node *BaseNode) GetLocation() string {
	if node.Tok == nil {
		return ""
	}

	return node.Tok.Locator.Locate()
}

func (node *BaseNode) SetParent(parent Node) {
	node.Parent = parent
}

func (node *BaseNode) GetParent() Node {
	return node.Parent
}

func (node *BaseNode) SetTyping(typing typing.Typing) {
	node.Typing = typing
}

func (node *BaseNode) GetTyping() typing.Typing {
	return node.Typing
}

func (node *BaseNode) SetScope(scope *symbolTable.Scope) {
	node.Scope = scope
}

func (node *BaseNode) GetScope() *symbolTable.Scope {
	return node.Scope
}

func (node *BaseNode) GetLocalScope() *symbolTable.Scope {
	localScope := node.Scope
	var localNode Node = node

	for localScope == nil && localNode.GetParent() != nil {
		localScope = localNode.GetParent().GetScope()
		localNode = localNode.GetParent()
	}

	return localScope
}

// Accept is part of visitor pattern.
func (node *BaseNode) Accept(visitor Visitor) {
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *BaseNode) VisitChildren(visitor Visitor) {
}
