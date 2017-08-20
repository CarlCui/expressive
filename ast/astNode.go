package ast

// Node represents a node in the ast
type Node interface {
	Accept(visitor Visitor)
	VisitChildren(visitor Visitor)
	SetParent(node Node)
	GetParent() Node
}
