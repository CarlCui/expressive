package ast

// BinaryOperatorNode represents a node with a binary operation (+, -, etc)
type BinaryOperatorNode struct {
	*BaseNode
	Lhs Node
	Rhs Node
}

// Accept is part of visitor pattern.
func (node *BinaryOperatorNode) Accept(visitor Visitor) {
	visitor.VisitEnterBinaryOepratorNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveBinaryOperatorNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *BinaryOperatorNode) VisitChildren(visitor Visitor) {
	node.Lhs.Accept(visitor)
	node.Rhs.Accept(visitor)
}
