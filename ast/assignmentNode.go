package ast

// AssignmentNode represents a node with assignment statement
type AssignmentNode struct {
	*BaseNode
	Identifier Node
	Expr       Node
}

// Accept is part of visitor pattern.
func (node *AssignmentNode) Accept(visitor Visitor) {
	visitor.VisitEnterAssignmentNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveAssignmentNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *AssignmentNode) VisitChildren(visitor Visitor) {
	node.Identifier.Accept(visitor)
	node.Expr.Accept(visitor)
}
