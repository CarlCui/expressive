package ast

// Visitor :see visitor pattern
type Visitor interface {
	VisitEnterProgramNode(node *ProgramNode)
	VisitLeaveProgramNode(node *ProgramNode)

	// stmts

	VisitEnterVariableDeclarationNode(node *VariableDeclarationNode)
	VisitLeaveVariableDeclarationNode(node *VariableDeclarationNode)

	VisitEnterAssignmentNode(node *AssignmentNode)
	VisitLeaveAssignmentNode(node *AssignmentNode)

	VisitEnterPrintNode(node *PrintNode)
	VisitLeavePrintNode(node *PrintNode)

	// exprs

	VisitEnterBinaryOepratorNode(node *BinaryOperatorNode)
	VisitLeaveBinaryOperatorNode(node *BinaryOperatorNode)

	VisitIntegerNode(node *IntegerNode)
	VisitFloatNode(node *FloatNode)
	VisitIdentifierNode(node *IdentifierNode)

	VisitErrorNode(node *ErrorNode)
}
