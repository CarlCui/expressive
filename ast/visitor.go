package ast

// Visitor :see visitor pattern
type Visitor interface {
	VisitEnterBinaryOepratorNode(node *BinaryOperatorNode)
	VisitLeaveBinaryOperatorNode(node *BinaryOperatorNode)

	VisitIntegerNode(node *IntegerNode)
	VisitFloatNode(node *FloatNode)
	VisitIdentifierNode(node *IdentifierNode)

	VisitErrorNode(node *ErrorNode)
}
