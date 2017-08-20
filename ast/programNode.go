package ast

import "github.com/carlcui/expressive/token"

// ProgramNode represents an integer constant node.
type ProgramNode struct {
	BaseNode
	Chilren []Node
}

// Accept is part of visitor pattern.
func (node *ProgramNode) Accept(visitor Visitor) {
	visitor.VisitEnterProgramNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveProgramNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *ProgramNode) VisitChildren(visitor Visitor) {
	for _, child := range node.Chilren {
		child.Accept(visitor)
	}
}

func (node *ProgramNode) Init(tok *token.Token) {
	node.BaseNode = CreateBaseNode(tok, nil)
}
