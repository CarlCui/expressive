package ast

import (
	"github.com/carlcui/expressive/symbolTable"
	"github.com/carlcui/expressive/typing"
)

// Node represents a node in the ast
type Node interface {
	Accept(visitor Visitor)
	VisitChildren(visitor Visitor)

	GetLocation() string

	SetParent(node Node)
	GetParent() Node

	SetScope(scope *symbolTable.Scope)
	GetScope() *symbolTable.Scope

	SetTyping(typing typing.Typing)
	GetTyping() typing.Typing
}

func Accept(node Node, visitor Visitor) {
	if node != nil {
		node.Accept(visitor)
	}
}
