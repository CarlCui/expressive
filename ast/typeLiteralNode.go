package ast

// TypeLiteralNode represents a node with a type literal
type TypeLiteralNode struct {
	*BaseNode
}

// Accept is part of visitor pattern.
func (node *TypeLiteralNode) Accept(visitor Visitor) {
	visitor.VisitTypeLiteralNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *TypeLiteralNode) VisitChildren(visitor Visitor) {
}
