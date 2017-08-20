package ast

// VariableDeclarationNode represents a node with variable declaration statement
type VariableDeclarationNode struct {
	BaseNode
	Identifier   Node
	DeclaredType Node
	Expr         Node
}

// Accept is part of visitor pattern.
func (node *VariableDeclarationNode) Accept(visitor Visitor) {
	visitor.VisitEnterVariableDeclarationNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveVariableDeclarationNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *VariableDeclarationNode) VisitChildren(visitor Visitor) {
	node.Identifier.Accept(visitor)
	node.DeclaredType.Accept(visitor)
	node.Expr.Accept(visitor)
}
