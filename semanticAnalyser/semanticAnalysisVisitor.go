package semanticAnalyser

import (
	"fmt"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/logger"
	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/symbolTable"
	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

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

func (visitor *SemanticAnalysisVisitor) VisitEnterBlockNode(node *ast.BlockNode) {
	localScope := node.GetLocalScope()
	newScope := symbolTable.CreateScope(localScope)
	node.SetScope(newScope)
}

func (visitor *SemanticAnalysisVisitor) VisitLeaveBlockNode(node *ast.BlockNode) {

}

// stmts

// VisitEnterVariableDeclarationNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterVariableDeclarationNode(node *ast.VariableDeclarationNode) {

}

// VisitLeaveVariableDeclarationNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeaveVariableDeclarationNode(node *ast.VariableDeclarationNode) {

	identifier := node.Identifier.(*ast.IdentifierNode)

	var resolvedTyping typing.Typing

	if node.DeclaredType == nil && node.Expr == nil {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.GetLocation(), "Missing variable type")
		return
	}

	if node.DeclaredType != nil {
		declaredTyping := node.DeclaredType.GetTyping()

		if node.Expr == nil {
			resolvedTyping = declaredTyping
		} else {
			exprTyping := node.Expr.GetTyping()

			if !exprTyping.Equals(declaredTyping) {
				visitor.log(node.GetLocation(),
					"variable declared as "+declaredTyping.String()+", "+
						"but expression evaluated to "+exprTyping.String())

				resolvedTyping = typing.ERROR_TYPE
			} else {
				resolvedTyping = exprTyping
			}
		}
	} else {
		resolvedTyping = node.Expr.GetTyping()
	}

	scope := node.GetLocalScope()

	if scope.VariableDeclared(identifier.Tok.Raw) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(identifier.GetLocation(), "variable \""+identifier.Tok.Raw+"\" has already been declared")
		return
	}

	if !scope.VariableCanBeShadowed(identifier.Tok.Raw) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(identifier.GetLocation(), "variable \""+identifier.Tok.Raw+"\" cannot be shadowed, thus already been declared")
		return
	}

	isDeclaringCanBeShadowed := func() bool {
		_, ok := node.Parent.(*ast.ForStmtNode)

		return ok
	}

	var binding *symbolTable.Binding

	if !isDeclaringCanBeShadowed() {
		binding = scope.CreateBinding(identifier.Tok.Raw, identifier.Tok.Locator, resolvedTyping)
	} else {
		binding = scope.CreateBindingCannotBeShadowed(identifier.Tok.Raw, identifier.Tok.Locator, resolvedTyping)
	}

	identifier.SetTyping(resolvedTyping)
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

	// check if assigning to constant
	if identifier, ok := node.LHS.(*ast.IdentifierNode); ok {
		binding := identifier.GetBinding()

		// If the binding is nil, then error is handled inside VisitIdentifierNode. If we pass a
		// null instance (not nil) to signature matching, it will log another error, which will
		// be miss-leading.
		if binding == nil {
			node.SetTyping(typing.ERROR_TYPE)
			return
		}

		if !binding.IsVariable {
			node.SetTyping(typing.ERROR_TYPE)
			visitor.log(identifier.GetLocation(), "variable cannot be re-assigned")
			return
		}
	}

	declaredType := node.LHS.GetTyping()
	exprType := node.RHS.GetTyping()

	if !declaredType.Equals(exprType) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.RHS.GetLocation(), "variable declared as "+declaredType.String()+", "+
			" but got "+exprType.String())
		return
	}

	// in the case of a compound assignment
	if node.Operator != signature.VOID_OPERATOR && !signature.HasSignature(node.Operator, declaredType, exprType) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.TypeCheckError(node.GetLocation(), node.Operator, declaredType, exprType)
		return
	}

	node.SetTyping(typing.VOID)
}

func (visitor *SemanticAnalysisVisitor) VisitEnterIncDecNode(node *ast.IncDecNode) {

}

func (visitor *SemanticAnalysisVisitor) VisitLeaveIncDecNode(node *ast.IncDecNode) {
	if _, ok := node.LHS.(*ast.IdentifierNode); !ok {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.LHS.GetLocation(), "left-hand side expression must be addressable.")
		return
	}

	lhsTyping := node.LHS.GetTyping()

	if !signature.HasSignature(signature.ADD, lhsTyping, typing.INT) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.TypeCheckError(node.GetLocation(), signature.ADD, lhsTyping, typing.INT)
		return
	}

	node.SetTyping(typing.VOID)
}

// VisitEnterPrintNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterPrintNode(node *ast.PrintNode) {

}

// VisitLeavePrintNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeavePrintNode(node *ast.PrintNode) {
	stringExprTyping := node.StringExpr.GetTyping()
	if !stringExprTyping.Equals(typing.STRING) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.GetLocation(), "requires string type, but got "+stringExprTyping.String())
		return
	}

	node.SetTyping(typing.VOID)
}

func (visitor *SemanticAnalysisVisitor) VisitEnterIfStmtNode(node *ast.IfStmtNode) {

}

func (visitor *SemanticAnalysisVisitor) VisitLeaveIfStmtNode(node *ast.IfStmtNode) {
	for _, conditionExpr := range node.ConditionExprs {
		conditionExprTyping := conditionExpr.GetTyping()
		if !conditionExprTyping.Equals(typing.BOOL) {
			node.SetTyping(typing.ERROR_TYPE)
			visitor.log(conditionExpr.GetLocation(), "requires boolean type, but got "+conditionExprTyping.String())
			return
		}
	}
}

func (visitor *SemanticAnalysisVisitor) VisitEnterWhileStmtNode(node *ast.WhileStmtNode) {

}

func (visitor *SemanticAnalysisVisitor) VisitLeaveWhileStmtNode(node *ast.WhileStmtNode) {
	conditionExprTyping := node.ConditionExpr.GetTyping()
	if !conditionExprTyping.Equals(typing.BOOL) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.ConditionExpr.GetLocation(), "requires boolean type, but got "+conditionExprTyping.String())
		return
	}
}

func (visitor *SemanticAnalysisVisitor) VisitEnterForStmtNode(node *ast.ForStmtNode) {
	localScope := node.GetLocalScope()
	newScope := symbolTable.CreateScope(localScope)
	node.SetScope(newScope)
}

func (visitor *SemanticAnalysisVisitor) VisitEnterForStmtNodeBeforeBlockNode(node *ast.ForStmtNode) {
	if node.ConditionExpr != nil {
		conditionExprTyping := node.ConditionExpr.GetTyping()
		if !conditionExprTyping.Equals(typing.BOOL) {
			node.SetTyping(typing.ERROR_TYPE)
			visitor.log(node.ConditionExpr.GetLocation(), "requires boolean type, but got "+conditionExprTyping.String())
			return
		}
	}
}

func (visitor *SemanticAnalysisVisitor) VisitLeaveForStmtNode(node *ast.ForStmtNode) {
	node.SetTyping(typing.VOID)
}

func (visitor *SemanticAnalysisVisitor) VisitEnterSwitchStmtNode(node *ast.SwitchStmtNode) {

}

func (visitor *SemanticAnalysisVisitor) VisitLeaveSwitchStmtNode(node *ast.SwitchStmtNode) {
	testExprTyping := node.TestExpr.GetTyping()

	for _, caseExpr := range node.CaseExprs {
		caseExprTyping := caseExpr.GetTyping()

		if !testExprTyping.Equals(caseExprTyping) {
			node.SetTyping(typing.ERROR_TYPE)
			visitor.log(caseExpr.GetLocation(), "has type "+caseExprTyping.String()+"does not match type of "+testExprTyping.String())
		}
	}

	node.SetTyping(typing.VOID)
}

func (visitor *SemanticAnalysisVisitor) VisitBreakNode(node *ast.BreakNode) {
	if node.FindNearestValidStatementNode() != nil {
		node.SetTyping(typing.VOID)
		return
	}

	node.SetTyping(typing.ERROR_TYPE)
	visitor.log(node.GetLocation(), "has to be inside of for, while or switch statement block")
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

	operator := node.Operator

	if !signature.HasSignature(operator, typing1, typing2, typing3) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.TypeCheckError(node.GetLocation(), operator, typing1, typing2, typing3)
		return
	}

	resultTyping := signature.ResultTyping(operator, typing1, typing2, typing3)

	node.SetTyping(resultTyping)
}

// VisitEnterBinaryOepratorNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterBinaryOepratorNode(node *ast.BinaryOperatorNode) {

}

// VisitLeaveBinaryOperatorNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeaveBinaryOperatorNode(node *ast.BinaryOperatorNode) {
	lhsTyping := node.Lhs.GetTyping()
	rhsTyping := node.Rhs.GetTyping()

	operator := node.Operator

	if !signature.HasSignature(operator, lhsTyping, rhsTyping) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.TypeCheckError(node.GetLocation(), operator, lhsTyping, rhsTyping)
		return
	}

	resultTyping := signature.ResultTyping(operator, lhsTyping, rhsTyping)

	node.SetTyping(resultTyping)
}

// VisitEnterUnaryOperatorNode do something
func (visitor *SemanticAnalysisVisitor) VisitEnterUnaryOperatorNode(node *ast.UnaryOperatorNode) {

}

// VisitLeaveUnaryOperatorNode do something
func (visitor *SemanticAnalysisVisitor) VisitLeaveUnaryOperatorNode(node *ast.UnaryOperatorNode) {
	paramTyping := node.Expr.GetTyping()

	operator := node.Operator

	if !signature.HasSignature(operator, paramTyping) {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.TypeCheckError(node.GetLocation(), operator, paramTyping)
		return
	}

	resultTyping := signature.ResultTyping(operator, paramTyping)

	node.SetTyping(resultTyping)
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

// VisitCharacterNode do something
func (visitor *SemanticAnalysisVisitor) VisitCharacterNode(node *ast.CharacterNode) {
	node.SetTyping(typing.CHAR)
}

// VisitStringNode do something
func (visitor *SemanticAnalysisVisitor) VisitStringNode(node *ast.StringNode) {
	node.SetTyping(typing.STRING)
}

// VisitIdentifierNode do something
func (visitor *SemanticAnalysisVisitor) VisitIdentifierNode(node *ast.IdentifierNode) {
	if node.IsBeingDeclared() {
		return
	}

	binding := node.FindVariableBinding()

	if binding == nil {
		node.SetTyping(typing.ERROR_TYPE)
		visitor.log(node.GetLocation(), "variable \""+node.Tok.Raw+"\" used before declared")
		return
	}

	node.SetTyping(binding.GetTyping())
	node.SetBinding(binding)
}

// VisitBooleanNode do something
func (visitor *SemanticAnalysisVisitor) VisitBooleanNode(node *ast.BooleanNode) {
	node.SetTyping(typing.BOOL)
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

func (visitor *SemanticAnalysisVisitor) TypeCheckError(location string, key interface{}, params ...typing.Typing) {
	err := fmt.Errorf("%v does not support operation on %v", key, params)
	visitor.log(location, err.Error())
}

func (visitor *SemanticAnalysisVisitor) log(location string, message string) {
	visitor.logger.Log(location, message)
}
