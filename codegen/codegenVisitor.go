package codegen

import (
	"fmt"
	"strconv"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/logger"
)

// CodegenVisitor visits each node and generates llvm IR.
type CodegenVisitor struct {
	logger                 logger.Logger
	labeller               *Labeller
	constants              *Fragment // global constants
	codeMap                map[ast.Node]*Fragment
	localIdentifierTracker *LocalIdentifierTracker
}

// Init with a logger
func (visitor *CodegenVisitor) Init(logger logger.Logger) {
	visitor.logger = logger
	visitor.labeller = &Labeller{0}
	visitor.constants = NewFragment(VOID, nil)
	visitor.codeMap = make(map[ast.Node]*Fragment)
	visitor.localIdentifierTracker = &LocalIdentifierTracker{0}
}

func (visitor *CodegenVisitor) newVoidCode(node ast.Node) *Fragment {
	if _, exists := visitor.codeMap[node]; exists {
		panic(fmt.Sprintf("Code for node %v already exists.", node))
	}

	fragment := NewFragment(VOID, visitor.localIdentifierTracker)
	visitor.codeMap[node] = fragment

	return fragment
}

func (visitor *CodegenVisitor) newValueCode(node ast.Node) *Fragment {
	if _, exists := visitor.codeMap[node]; exists {
		panic(fmt.Sprintf("Code for node %v already exists.", node))
	}

	fragment := NewFragment(VALUE, visitor.localIdentifierTracker)
	visitor.codeMap[node] = fragment

	return fragment
}

func (visitor *CodegenVisitor) newAddressCode(node ast.Node) *Fragment {
	if _, exists := visitor.codeMap[node]; exists {
		panic(fmt.Sprintf("Code for node %v already exists.", node))
	}

	fragment := NewFragment(ADDRESS, visitor.localIdentifierTracker)
	visitor.codeMap[node] = fragment

	return fragment
}

func (visitor *CodegenVisitor) getAndRemoveCode(node ast.Node) *Fragment {
	fragment, exists := visitor.codeMap[node]

	if !exists {
		panic(fmt.Sprintf("Code for node %v does not exist.", node))
	}

	delete(visitor.codeMap, node)

	return fragment
}

func (visitor *CodegenVisitor) removeVoidCode(node ast.Node) *Fragment {
	fragment := visitor.getAndRemoveCode(node)

	if fragment.ResultType != VOID {
		panic(fmt.Sprintf("Code fragment does not produce void result: %v", node))
	}

	return fragment
}

func (visitor *CodegenVisitor) removeAddressCode(node ast.Node) *Fragment {
	fragment := visitor.getAndRemoveCode(node)

	if fragment.ResultType != ADDRESS {
		panic(fmt.Sprintf("Code fragment does not produce address result: %v", node))
	}

	return fragment
}

func (visitor *CodegenVisitor) removeValueCode(node ast.Node) *Fragment {
	fragment := visitor.getAndRemoveCode(node)

	if fragment.ResultType != VALUE {
		panic(fmt.Sprintf("Code fragment does not produce value result: %v", node))
	}

	if fragment.ResultType == ADDRESS {
		visitor.turnAddressIntoValue(node, fragment)
	}

	return fragment
}

func (visitor *CodegenVisitor) turnAddressIntoValue(node ast.Node, fragment *Fragment) {
	fragment.ResultType = VALUE

	result := fragment.result
	typing := node.GetTyping()
	irType := typing.IrType()

	fragment.AddOperation("load %v, %v* %v", irType, irType, result)
}

// VisitEnterProgramNode creates program scope
func (visitor *CodegenVisitor) VisitEnterProgramNode(node *ast.ProgramNode) {

}

// VisitLeaveProgramNode closes program scope
func (visitor *CodegenVisitor) VisitLeaveProgramNode(node *ast.ProgramNode) {
	// generates main function
	fragment := visitor.newVoidCode(node)

	fragment.AddInstruction("define i32 @main() {")

	for _, child := range node.Chilren {
		fragment.Append(visitor.removeVoidCode(child))
	}

	fragment.AddInstruction("}")
}

// stmts

// VisitEnterVariableDeclarationNode do something
func (visitor *CodegenVisitor) VisitEnterVariableDeclarationNode(node *ast.VariableDeclarationNode) {

}

// VisitLeaveVariableDeclarationNode do something
func (visitor *CodegenVisitor) VisitLeaveVariableDeclarationNode(node *ast.VariableDeclarationNode) {
	fragment := visitor.newVoidCode(node)

	identifierNode := node.Identifier.(*ast.IdentifierNode)
	typing := identifierNode.GetTyping()
	irType := typing.IrType()
	alignment := typing.Size()

	variable := AsLocalVariable(identifierNode.Tok.Raw)

	// allocate space
	fragment.AddInstruction("%v = alloca %v, align %v", variable, irType, alignment)

	if node.Expr == nil {
		// load default value
		fragment.AddInstruction("store %v %v, %v* %v, align %v", irType, 0, irType, variable, alignment)
	} else {
		exprFragment := visitor.removeValueCode(node.Expr)

		exprResultVariable := exprFragment.GetResult()

		fragment.Append(exprFragment)

		fragment.AddInstruction("store %v %v, %v* %v, align %v", irType, exprResultVariable, irType, variable, alignment)
	}
}

// VisitEnterAssignmentNode do something
func (visitor *CodegenVisitor) VisitEnterAssignmentNode(node *ast.AssignmentNode) {
}

// VisitLeaveAssignmentNode do something
func (visitor *CodegenVisitor) VisitLeaveAssignmentNode(node *ast.AssignmentNode) {
	fragment := visitor.newVoidCode(node)

	identifierNode := node.Identifier.(*ast.IdentifierNode)
	typing := identifierNode.GetTyping()
	irType := typing.IrType()
	alignment := typing.Size()

	variable := AsLocalVariable(identifierNode.Tok.Raw)

	exprFragment := visitor.removeValueCode(node.Expr)

	exprResultVariable := exprFragment.GetResult()

	fragment.Append(exprFragment)

	fragment.AddInstruction("store %v %v, %v* %v, align %v", irType, exprResultVariable, irType, variable, alignment)
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
	fragment := visitor.newValueCode(node)

	fragment1 := visitor.removeValueCode(node.Expr1)
	fragment2 := visitor.removeValueCode(node.Expr2)
	fragment3 := visitor.removeValueCode(node.Expr3)

	operator := node.Operator
	typing := node.GetTyping()

	operatorCodegen := NewOperatorCodegen(fragment, operator, typing, visitor.labeller, fragment1, fragment2, fragment3)

	operatorCodegen.GenerateCode()
}

// VisitEnterBinaryOepratorNode do something
func (visitor *CodegenVisitor) VisitEnterBinaryOepratorNode(node *ast.BinaryOperatorNode) {

}

// VisitLeaveBinaryOperatorNode do something
func (visitor *CodegenVisitor) VisitLeaveBinaryOperatorNode(node *ast.BinaryOperatorNode) {
	fragment := visitor.newValueCode(node)

	fragment1 := visitor.removeValueCode(node.Lhs)
	fragment2 := visitor.removeValueCode(node.Rhs)

	operator := node.Operator
	typing := node.GetTyping()

	operatorCodegen := NewOperatorCodegen(fragment, operator, typing, visitor.labeller, fragment1, fragment2)

	operatorCodegen.GenerateCode()
}

// VisitEnterUnaryOperatorNode do something
func (visitor *CodegenVisitor) VisitEnterUnaryOperatorNode(node *ast.UnaryOperatorNode) {

}

// VisitLeaveUnaryOperatorNode do something
func (visitor *CodegenVisitor) VisitLeaveUnaryOperatorNode(node *ast.UnaryOperatorNode) {
	fragment := visitor.newValueCode(node)

	fragment1 := visitor.removeValueCode(node.Expr)

	operator := node.Operator
	typing := node.GetTyping()

	operatorCodegen := NewOperatorCodegen(fragment, operator, typing, visitor.labeller, fragment1)

	operatorCodegen.GenerateCode()
}

// literal nodes

// VisitIntegerNode do something
func (visitor *CodegenVisitor) VisitIntegerNode(node *ast.IntegerNode) {
	fragment := visitor.newValueCode(node)

	fragment.result = strconv.Itoa(node.Val)
}

// VisitFloatNode do something
func (visitor *CodegenVisitor) VisitFloatNode(node *ast.FloatNode) {
	fragment := visitor.newValueCode(node)

	fragment.result = strconv.FormatFloat(float64(node.Val), 'f', 6, 32)
}

// VisitCharacterNode do something
func (visitor *CodegenVisitor) VisitCharacterNode(node *ast.CharacterNode) {

}

// VisitStringNode do something
func (visitor *CodegenVisitor) VisitStringNode(node *ast.StringNode) {

}

// VisitIdentifierNode do something
func (visitor *CodegenVisitor) VisitIdentifierNode(node *ast.IdentifierNode) {
	fragment := visitor.newAddressCode(node)

	identifier := node.Tok.Raw
	typing := node.GetTyping()
	irType := typing.IrType()
	alignment := typing.Size()

	fragment.AddInstruction("%%v = alloca %v, align %v", identifier, irType, alignment)

	fragment.result = AsLocalVariable(identifier)
}

// VisitBooleanNode do something
func (visitor *CodegenVisitor) VisitBooleanNode(node *ast.BooleanNode) {
	fragment := visitor.newValueCode(node)

	var booleanValue string

	if node.Val == true {
		booleanValue = "1"
	} else {
		booleanValue = "0"
	}

	fragment.result = booleanValue
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