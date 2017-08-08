package ast

// IdentifierNode represents an identifier node.
type IdentifierNode struct {
	BaseNode
}

// Accept is part of visitor pattern.
func (node *IdentifierNode) Accept(visitor Visitor) {
	visitor.VisitIdentifierNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *IdentifierNode) VisitChildren(visitor Visitor) {

}
