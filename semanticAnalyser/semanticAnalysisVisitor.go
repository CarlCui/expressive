package semanticAnalyser

import "github.com/carlcui/expressive/ast"

// SemanticAnalysisVisitor is the general semantic analyser using visitor pattern
type SemanticAnalysisVisitor int

// VisitEnterProgramNode creates program scope
func (visitor SemanticAnalysisVisitor) VisitEnterProgramNode(node *ast.ProgramNode) {

}

// VisitLeaveProgramNode closes program scope
func (visitor SemanticAnalysisVisitor) VisitLeaveProgramNode(node *ast.ProgramNode) {

}

// stmts

// VisitEnterVariableDeclarationNode do something
func (visitor SemanticAnalysisVisitor) VisitEnterVariableDeclarationNode(node *ast.VariableDeclarationNode) {

}

// VisitLeaveVariableDeclarationNode do something
func (visitor SemanticAnalysisVisitor) VisitLeaveVariableDeclarationNode(node *ast.VariableDeclarationNode) {

}

// VisitEnterAssignmentNode do something
func (visitor SemanticAnalysisVisitor) VisitEnterAssignmentNode(node *ast.AssignmentNode) {

}

// VisitLeaveAssignmentNode do something
func (visitor SemanticAnalysisVisitor) VisitLeaveAssignmentNode(node *ast.AssignmentNode) {

}

// VisitEnterPrintNode do something
func (visitor SemanticAnalysisVisitor) VisitEnterPrintNode(node *ast.PrintNode) {

}

// VisitLeavePrintNode do something
func (visitor SemanticAnalysisVisitor) VisitLeavePrintNode(node *ast.PrintNode) {

}

// exprs

// VisitEnterTernaryOperatorNode do something
func (visitor SemanticAnalysisVisitor) VisitEnterTernaryOperatorNode(node *ast.TernaryOperatorNode) {

}

// VisitLeaveTernaryOperatorNode do something
func (visitor SemanticAnalysisVisitor) VisitLeaveTernaryOperatorNode(node *ast.TernaryOperatorNode) {

}

// VisitEnterBinaryOepratorNode do something
func (visitor SemanticAnalysisVisitor) VisitEnterBinaryOepratorNode(node *ast.BinaryOperatorNode) {

}

// VisitLeaveBinaryOperatorNode do something
func (visitor SemanticAnalysisVisitor) VisitLeaveBinaryOperatorNode(node *ast.BinaryOperatorNode) {

}

// VisitEnterUnaryOperatorNode do something
func (visitor SemanticAnalysisVisitor) VisitEnterUnaryOperatorNode(node *ast.UnaryOperatorNode) {

}

// VisitLeaveUnaryOperatorNode do something
func (visitor SemanticAnalysisVisitor) VisitLeaveUnaryOperatorNode(node *ast.UnaryOperatorNode) {

}

// literal nodes

// VisitIntegerNode do something
func (visitor SemanticAnalysisVisitor) VisitIntegerNode(node *ast.IntegerNode) {

}

// VisitFloatNode do something
func (visitor SemanticAnalysisVisitor) VisitFloatNode(node *ast.FloatNode) {

}

// VisitIdentifierNode do something
func (visitor SemanticAnalysisVisitor) VisitIdentifierNode(node *ast.IdentifierNode) {

}

// VisitTypeLiteralNode do something
func (visitor SemanticAnalysisVisitor) VisitTypeLiteralNode(node *ast.TypeLiteralNode) {

}

// VisitErrorNode do something
func (visitor SemanticAnalysisVisitor) VisitErrorNode(node *ast.ErrorNode) {

}
