package codegen

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/logger"
)

type CodegenVisitor struct {
	logger logger.Logger
}

// VisitEnterProgramNode creates program scope
func (visitor *CodegenVisitor) VisitEnterProgramNode(node *ast.ProgramNode) {

}

// VisitLeaveProgramNode closes program scope
func (visitor *CodegenVisitor) VisitLeaveProgramNode(node *ast.ProgramNode) {

}

// stmts

// VisitEnterVariableDeclarationNode do something
func (visitor *CodegenVisitor) VisitEnterVariableDeclarationNode(node *ast.VariableDeclarationNode) {

}

// VisitLeaveVariableDeclarationNode do something
func (visitor *CodegenVisitor) VisitLeaveVariableDeclarationNode(node *ast.VariableDeclarationNode) {

}

// VisitEnterAssignmentNode do something
func (visitor *CodegenVisitor) VisitEnterAssignmentNode(node *ast.AssignmentNode) {
}

// VisitLeaveAssignmentNode do something
func (visitor *CodegenVisitor) VisitLeaveAssignmentNode(node *ast.AssignmentNode) {

}

// VisitEnterPrintNode do something
func (visitor *CodegenVisitor) VisitEnterPrintNode(node *ast.PrintNode) {

}

// VisitLeavePrintNode do something
func (visitor *CodegenVisitor) VisitLeavePrintNode(node *ast.PrintNode) {

}

// exprs

// VisitEnterTernaryOperatorNode do something
func (visitor *CodegenVisitor) VisitEnterTernaryOperatorNode(node *ast.TernaryOperatorNode) {

}

// VisitLeaveTernaryOperatorNode do something
func (visitor *CodegenVisitor) VisitLeaveTernaryOperatorNode(node *ast.TernaryOperatorNode) {

}

// VisitEnterBinaryOepratorNode do something
func (visitor *CodegenVisitor) VisitEnterBinaryOepratorNode(node *ast.BinaryOperatorNode) {

}

// VisitLeaveBinaryOperatorNode do something
func (visitor *CodegenVisitor) VisitLeaveBinaryOperatorNode(node *ast.BinaryOperatorNode) {

}

// VisitEnterUnaryOperatorNode do something
func (visitor *CodegenVisitor) VisitEnterUnaryOperatorNode(node *ast.UnaryOperatorNode) {

}

// VisitLeaveUnaryOperatorNode do something
func (visitor *CodegenVisitor) VisitLeaveUnaryOperatorNode(node *ast.UnaryOperatorNode) {

}

// literal nodes

// VisitIntegerNode do something
func (visitor *CodegenVisitor) VisitIntegerNode(node *ast.IntegerNode) {

}

// VisitFloatNode do something
func (visitor *CodegenVisitor) VisitFloatNode(node *ast.FloatNode) {

}

// VisitCharacterNode do something
func (visitor *CodegenVisitor) VisitCharacterNode(node *ast.CharacterNode) {

}

// VisitStringNode do something
func (visitor *CodegenVisitor) VisitStringNode(node *ast.StringNode) {

}

// VisitIdentifierNode do something
func (visitor *CodegenVisitor) VisitIdentifierNode(node *ast.IdentifierNode) {

}

// VisitBooleanNode do something
func (visitor *CodegenVisitor) VisitBooleanNode(node *ast.BooleanNode) {

}

// VisitTypeLiteralNode do something
func (visitor *CodegenVisitor) VisitTypeLiteralNode(node *ast.TypeLiteralNode) {

}

// VisitErrorNode should not happen during codegen
func (visitor *CodegenVisitor) VisitErrorNode(node *ast.ErrorNode) {
	panic(node.GetLocation() + ": unexpected error node")
}

func (visitor *CodegenVisitor) log(location string, message string) {
	visitor.logger.Log(location, message)
}
