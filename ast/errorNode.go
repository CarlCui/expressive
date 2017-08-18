package ast

// ErrorNode represents a node with syntax error.
type ErrorNode struct {
	BaseNode
}

// Accept is part of visitor pattern.
func (node *ErrorNode) Accept(visitor Visitor) {
	visitor.VisitErrorNode(node)
}

// VisitChildren is part of visitor pattern. Error node does not have any child.
func (node *ErrorNode) VisitChildren(visitor Visitor) {

}
