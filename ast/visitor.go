package ast

// Visitor :see visitor pattern
type Visitor interface {
	VisitEnterProgramNode(node *ProgramNode)
	VisitLeaveProgramNode(node *ProgramNode)

	VisitEnterBlockNode(node *BlockNode)
	VisitLeaveBlockNode(node *BlockNode)

	// stmts

	VisitEnterVariableDeclarationNode(node *VariableDeclarationNode)
	VisitLeaveVariableDeclarationNode(node *VariableDeclarationNode)

	VisitEnterAssignmentNode(node *AssignmentNode)
	VisitLeaveAssignmentNode(node *AssignmentNode)

	VisitEnterIfStmtNode(node *IfStmtNode)
	VisitLeaveIfStmtNode(node *IfStmtNode)

	VisitEnterForStmtNode(node *ForStmtNode)
	VisitEnterForStmtNodeBeforeBlockNode(node *ForStmtNode)
	VisitLeaveForStmtNode(node *ForStmtNode)

	VisitEnterPrintNode(node *PrintNode)
	VisitLeavePrintNode(node *PrintNode)

	// exprs

	VisitEnterTernaryOperatorNode(node *TernaryOperatorNode)
	VisitLeaveTernaryOperatorNode(node *TernaryOperatorNode)

	VisitEnterBinaryOepratorNode(node *BinaryOperatorNode)
	VisitLeaveBinaryOperatorNode(node *BinaryOperatorNode)

	VisitEnterUnaryOperatorNode(node *UnaryOperatorNode)
	VisitLeaveUnaryOperatorNode(node *UnaryOperatorNode)

	// literal nodes

	VisitIntegerNode(node *IntegerNode)
	VisitFloatNode(node *FloatNode)
	VisitBooleanNode(node *BooleanNode)
	VisitCharacterNode(node *CharacterNode)
	VisitStringNode(node *StringNode)
	VisitIdentifierNode(node *IdentifierNode)

	VisitTypeLiteralNode(node *TypeLiteralNode)

	VisitErrorNode(node *ErrorNode)
}
