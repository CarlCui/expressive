package semanticAnalyser

import "github.com/carlcui/expressive/ast"
import "github.com/carlcui/expressive/symbolTable"
import "github.com/carlcui/expressive/typing"
import "github.com/carlcui/expressive/token"
import "github.com/carlcui/expressive/logger"

// SemanticAnalysisVisitor is the general semantic analyser using visitor pattern
type SemanticAnalysisVisitor struct {
	logger logger.Logger
}

// VisitEnterProgramNode creates program scope
func (visitor *SemanticAnalysisVisitor) VisitEnterProgramNode(node *ast.ProgramNode) {
	scope := symbolTable.CreateScope(nil)
	node.SetScope(scope)
}

// VisitLeaveProgramNode closes program scope
func (visitor *SemanticAnalysisVisitor) VisitLeaveProgramNode(node *ast.ProgramNode) {

}

// stmts

// VisitEnterVariableDeclarationNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterVariableDeclarationNode(node *ast.VariableDeclarationNode) {

}

// VisitLeaveVariableDeclarationNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeaveVariableDeclarationNode(node *ast.VariableDeclarationNode) {

	identifier := node.Identifier.(*ast.IdentifierNode)

	var resolvedTyping typing.Typing

	declaredTyping := node.DeclaredType.GetTyping()

	if node.DeclaredType == nil && node.Expr == nil {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.GetLocation(), "Missing variable type")
		return
	}

	if node.DeclaredType != nil {
		exprTyping := node.Expr.GetTyping()

		if !exprTyping.Equals(declaredTyping) {
			visitor.log(node.GetLocation(),
				"variable declared as "+declaredTyping.String()+","+
					"but expression evaluated to "+exprTyping.String())

			resolvedTyping = typing.ERROR_TYPE
		} else {
			resolvedTyping = exprTyping
		}
	} else {
		resolvedTyping = node.Expr.GetTyping()
	}

	scope := node.GetLocalScope()

	if scope.VariableDeclared(identifier.Tok.Raw) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(identifier.GetLocation(), "variable "+identifier.Tok.Raw+" already declared at ")
		return
	}

	binding := scope.CreateBinding(identifier.Tok.Raw, identifier.Tok.Locator, resolvedTyping)
	identifier.SetBinding(binding)

	if node.Tok.TokenType == token.CONST {
		binding.IsVariable = false
	}

	node.SetTyping(typing.VOID)
}

// VisitEnterAssignmentNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterAssignmentNode(node *ast.AssignmentNode) {
}

// VisitLeaveAssignmentNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeaveAssignmentNode(node *ast.AssignmentNode) {

	identifier := node.Identifier.(*ast.IdentifierNode)

	binding := identifier.GetBinding()

	if !binding.IsVariable {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(identifier.GetLocation(), "variable cannot be re-assigned")
		return
	}

	declaredType := identifier.GetTyping()
	exprType := node.Expr.GetTyping()

	if !declaredType.Equals(exprType) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.Expr.GetLocation(), "variable declared as "+declaredType.String()+", "+
			"but got "+exprType.String())
		return
	}

	node.SetTyping(typing.VOID)
}

// VisitEnterPrintNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterPrintNode(node *ast.PrintNode) {

}

// VisitLeavePrintNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeavePrintNode(node *ast.PrintNode) {
	exprTyping := node.Expr.GetTyping()
	if !exprTyping.Equals(typing.STRING) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.GetLocation(), "requires string type, but got "+exprTyping.String())
		return
	}

	node.SetTyping(typing.VOID)
}

// exprs

// VisitEnterTernaryOperatorNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterTernaryOperatorNode(node *ast.TernaryOperatorNode) {

}

// VisitLeaveTernaryOperatorNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeaveTernaryOperatorNode(node *ast.TernaryOperatorNode) {
	typing1 := node.Expr1.GetTyping()
	typing2 := node.Expr2.GetTyping()
	typing3 := node.Expr3.GetTyping()

	if !typing1.Equals(typing.BOOL) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.Expr1.GetLocation(), "requires bool, but got "+typing1.String())
		return
	}

	if !typing2.Equals(typing3) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.Expr2.GetLocation(),
			"ternary operator ? : requires both expression to evaluate to the same type, but got "+typing2.String()+" and "+typing3.String())
		return
	}

	node.SetTyping(typing2)
}

// VisitEnterBinaryOepratorNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterBinaryOepratorNode(node *ast.BinaryOperatorNode) {

}

// VisitLeaveBinaryOperatorNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeaveBinaryOperatorNode(node *ast.BinaryOperatorNode) {
	lhsTyping := node.Lhs.GetTyping()
	rhsTyping := node.Rhs.GetTyping()

	if !lhsTyping.Equals(rhsTyping) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.Lhs.GetLocation(), "operator "+node.Tok.Raw+" does not support "+
			lhsTyping.String()+" and "+rhsTyping.String())
		return
	}

	// get a way to do this efficiently
}

// VisitEnterUnaryOperatorNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterUnaryOperatorNode(node *ast.UnaryOperatorNode) {

}

// VisitLeaveUnaryOperatorNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeaveUnaryOperatorNode(node *ast.UnaryOperatorNode) {

}

// literal nodes

// VisitIntegerNode do something
func (visitor *SemanticAnalysisVisitor) VisitIntegerNode(node *ast.IntegerNode) {
	node.SetTyping(typing.INT)
}

// VisitFloatNode do something
func (visitor *SemanticAnalysisVisitor) VisitFloatNode(node *ast.FloatNode) {
	node.SetTyping(typing.FLOAT)
}

// VisitIdentifierNode do something
func (visitor *SemanticAnalysisVisitor) VisitIdentifierNode(node *ast.IdentifierNode) {
	if !node.IsBeingDeclared() {
		binding := node.FindVariableBinding()

		if binding == nil {
			node.SetTyping(typing.ERROR_TYPE)
			visitor.log(node.GetLocation(), "variable "+node.Tok.Raw+" used before declared")
		}

		node.SetTyping(binding.GetTyping())
		node.SetBinding(binding)
	}
}

// VisitTypeLiteralNode do something
func (visitor *SemanticAnalysisVisitor) VisitTypeLiteralNode(node *ast.TypeLiteralNode) {
	switch node.Tok.TokenType {
	case token.INT_KEYWORD:
		node.SetTyping(typing.INT)
		break
	case token.FLOAT_KEYWORD:
		node.SetTyping(typing.FLOAT)
		break
	case token.CHAR_KEYWORD:
		node.SetTyping(typing.CHAR)
		break
	case token.BOOL_KEYWORD:
		node.SetTyping(typing.BOOL)
		break
	case token.STRING_KEYWORD:
		node.SetTyping(typing.STRING)
		break
	default:
		node.SetTyping(typing.NO_TYPE)
	}
}

// VisitErrorNode do something
func (visitor *SemanticAnalysisVisitor) VisitErrorNode(node *ast.ErrorNode) {
	node.SetTyping(typing.ERROR_TYPE)
}

func (visitor *SemanticAnalysisVisitor) log(location string, message string) {
	visitor.logger.Log(location, message)
}
