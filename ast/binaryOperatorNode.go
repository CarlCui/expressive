package ast

// BinaryOperatorNode represents a node with a binary operation (+, -, etc)
type BinaryOperatorNode struct {
	BaseNode
	lhs Node
	rhs Node
}

// Accept is part of visitor pattern.
func (node *BinaryOperatorNode) Accept(visitor Visitor) {
	visitor.VisitEnterBinaryOepratorNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveBinaryOperatorNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *BinaryOperatorNode) VisitChildren(visitor Visitor) {
	node.lhs.Accept(visitor)
	node.rhs.Accept(visitor)
}
