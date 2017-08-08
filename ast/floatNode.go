package ast

// FloatNode represents a float constant node.
type FloatNode struct {
	BaseNode
	val float32
}

// Accept is part of visitor pattern.
func (node *FloatNode) Accept(visitor Visitor) {
	visitor.VisitFloatNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *FloatNode) VisitChildren(visitor Visitor) {

}
