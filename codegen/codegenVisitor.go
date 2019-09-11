package codegen

import (
	"fmt"
	"strconv"

	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/typing"

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
	externals              *Fragment
}

func (visitor *CodegenVisitor) externalFragment() *Fragment {
	fragment := NewFragment(VOID, nil)

	fragment.AddInstruction("declare i32 @printf(i8* noalias nocapture, ...) nounwind")

	return fragment
}

// Init with a logger
func (visitor *CodegenVisitor) Init(logger logger.Logger) {
	visitor.logger = logger
	visitor.labeller = &Labeller{0}
	visitor.constants = NewFragment(VOID, &GlobalIdentifierTracker{0})
	visitor.codeMap = make(map[ast.Node]*Fragment)
	visitor.localIdentifierTracker = &LocalIdentifierTracker{0}
	visitor.externals = visitor.externalFragment()
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

	if fragment.ResultType != VALUE && fragment.ResultType != ADDRESS {
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

	fragment.AddInstruction("ret i32 0")
	fragment.AddInstruction("}")
}

func (visitor *CodegenVisitor) VisitEnterBlockNode(node *ast.BlockNode) {

}

func (visitor *CodegenVisitor) VisitLeaveBlockNode(node *ast.BlockNode) {
	fragment := visitor.newVoidCode(node)

	for _, child := range node.Stmts {
		fragment.Append(visitor.removeVoidCode(child))
	}
}

// stmts

// VisitEnterVariableDeclarationNode do something
func (visitor *CodegenVisitor) VisitEnterVariableDeclarationNode(node *ast.VariableDeclarationNode) {

}

// VisitLeaveVariableDeclarationNode do something
func (visitor *CodegenVisitor) VisitLeaveVariableDeclarationNode(node *ast.VariableDeclarationNode) {
	fragment := visitor.newVoidCode(node)

	identifierNode := node.Identifier.(*ast.IdentifierNode)
	identifierTyping := identifierNode.GetTyping()
	irType := identifierTyping.IrType()
	alignment := identifierTyping.Size()

	variable := AsLocalVariable(identifierNode.LocalIdentifier())

	// allocate space
	fragment.AddInstruction("%v = alloca %v, align %v", variable, irType, alignment)

	if node.Expr == nil {
		// load default value
		defaultValue := "0"
		if identifierTyping == typing.FLOAT {
			defaultValue = "0.0"
		}

		fragment.AddInstruction("store %v %v, %v* %v, align %v", irType, defaultValue, irType, variable, alignment)
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

	variable := AsLocalVariable(identifierNode.LocalIdentifier())

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
	fragment := visitor.newVoidCode(node)

	stringExprFrag := visitor.removeAddressCode(node.StringExpr)

	instructionArgs := make([]interface{}, 0)

	localIdentifier := stringExprFrag.GetResult()

	instructionArgs = append(instructionArgs, localIdentifier)

	fragment.Append(stringExprFrag)

	callInstruction := "call i32 (i8*, ...) @printf(i8* %v"

	for _, arg := range node.Args {
		argFrag := visitor.removeValueCode(arg)
		localIdentifier = argFrag.GetResult()

		irType := arg.GetTyping().IrType()

		fragment.Append(argFrag)

		instructionArgs = append(instructionArgs, irType, localIdentifier)

		callInstruction += ", %v %v"
	}

	callInstruction += ")"

	fragment.AddInstruction(callInstruction, instructionArgs...)
}

func (visitor *CodegenVisitor) VisitEnterIfStmtNode(node *ast.IfStmtNode) {

}

func (visitor *CodegenVisitor) VisitLeaveIfStmtNode(node *ast.IfStmtNode) {
	fragment := visitor.newVoidCode(node)

	ifEndLabel := visitor.labeller.NewSet("if", "end")
	ifElseLabel := visitor.labeller.Label("if", "else")

	numberOfConditions := len(node.ConditionExprs)

	ifConditionLabels := make([]string, numberOfConditions)

	for i := range ifConditionLabels {
		ifConditionLabels[i] = visitor.labeller.Label("if", "condition", strconv.Itoa(i))
	}

	for i, conditionExpr := range node.ConditionExprs {
		conditionLabel := ifConditionLabels[i]
		blockLabel := visitor.labeller.Label("if", "block", strconv.Itoa(i))

		exprFragment := visitor.removeValueCode(conditionExpr)

		fragment.AddInstruction("br label %v", AsLocalVariable(conditionLabel))
		fragment.AddLabel(conditionLabel)

		fragment.Append(exprFragment)

		exprResult := exprFragment.GetResult()

		conditionResult := fragment.AddOperation("icmp eq i1 %v, 1", exprResult)

		finalLabel := ifEndLabel

		if i+1 != numberOfConditions {
			finalLabel = ifConditionLabels[i+1]
		} else if node.ElseBlock != nil {
			finalLabel = ifElseLabel
		}

		fragment.AddInstruction("br i1 %v, label %v, label %v", conditionResult, AsLocalVariable(blockLabel), AsLocalVariable(finalLabel))

		fragment.AddLabel(blockLabel)
		fragment.Append(visitor.removeVoidCode(node.ConditionBlocks[i]))
		fragment.AddInstruction("br label %v", AsLocalVariable(ifEndLabel))
	}

	if node.ElseBlock != nil {
		fragment.AddInstruction("br label %v", AsLocalVariable(ifElseLabel))
		fragment.AddLabel(ifElseLabel)
		fragment.Append(visitor.removeVoidCode(node.ElseBlock))
	}

	fragment.AddInstruction("br label %v", AsLocalVariable(ifEndLabel))
	fragment.AddLabel(ifEndLabel)
}

func (visitor *CodegenVisitor) VisitEnterWhileStmtNode(node *ast.WhileStmtNode) {
	node.EndLabel = visitor.labeller.NewSet("while", "end")
}

func (visitor *CodegenVisitor) VisitLeaveWhileStmtNode(node *ast.WhileStmtNode) {
	fragment := visitor.newVoidCode(node)

	whileStartLabel := visitor.labeller.NewSet("while", "start")
	whileEndLabel := node.EndLabel
	whileConditionExprLabel := visitor.labeller.Label("while", "condition")
	whileBlockLabel := visitor.labeller.Label("while", "block")

	fragment.AddInstruction("br label %v", AsLocalVariable(whileStartLabel))
	fragment.AddLabel(whileStartLabel)
	fragment.AddInstruction("br label %v", AsLocalVariable(whileConditionExprLabel))
	fragment.AddLabel(whileConditionExprLabel)

	conditionExprFragment := visitor.removeValueCode(node.ConditionExpr)

	fragment.Append(conditionExprFragment)

	conditionResult := fragment.AddOperation("icmp eq i1 %v, 1", conditionExprFragment.GetResult())

	fragment.AddInstruction("br i1 %v, label %v, label %v", conditionResult, AsLocalVariable(whileBlockLabel), AsLocalVariable(whileEndLabel))

	fragment.AddLabel(whileBlockLabel)
	fragment.Append(visitor.removeVoidCode(node.Block))

	fragment.AddInstruction("br label %v", AsLocalVariable(whileConditionExprLabel))

	fragment.AddLabel(whileEndLabel)
}

func (visitor *CodegenVisitor) VisitEnterForStmtNode(node *ast.ForStmtNode) {
	node.EndLabel = visitor.labeller.NewSet("for", "end")
}

func (visitor *CodegenVisitor) VisitEnterForStmtNodeBeforeBlockNode(node *ast.ForStmtNode) {

}

func (visitor *CodegenVisitor) VisitLeaveForStmtNode(node *ast.ForStmtNode) {
	fragment := visitor.newVoidCode(node)

	forStartLabel := visitor.labeller.NewSet("for", "start")
	forEndLabel := node.EndLabel

	initializationLabel := visitor.labeller.Label("for", "initialization")
	conditionExprLabel := visitor.labeller.Label("for", "condition")
	iterationStmtLabel := visitor.labeller.Label("for", "iteration")
	blockLabel := visitor.labeller.Label("for", "block")

	fragment.AddInstruction("br label %v", AsLocalVariable(forStartLabel))
	fragment.AddLabel(forStartLabel)

	fragment.AddInstruction("br label %v", AsLocalVariable(initializationLabel))
	fragment.AddLabel(initializationLabel)
	if node.InitializationStmt != nil {
		fragment.Append(visitor.removeVoidCode(node.InitializationStmt))
	}

	fragment.AddInstruction("br label %v", AsLocalVariable(conditionExprLabel))
	fragment.AddLabel(conditionExprLabel)

	if node.ConditionExpr != nil {
		conditionExprFragment := visitor.removeValueCode(node.ConditionExpr)

		fragment.Append(conditionExprFragment)

		conditionResult := fragment.AddOperation("icmp eq i1 %v, 1", conditionExprFragment.GetResult())

		fragment.AddInstruction("br i1 %v, label %v, label %v", conditionResult, AsLocalVariable(blockLabel), AsLocalVariable(forEndLabel))
	}

	fragment.AddLabel(blockLabel)
	fragment.Append(visitor.removeVoidCode(node.Block))

	fragment.AddInstruction("br label %v", AsLocalVariable(iterationStmtLabel))
	fragment.AddLabel(iterationStmtLabel)
	if node.IterationStmt != nil {
		fragment.Append(visitor.removeVoidCode(node.IterationStmt))
	}

	fragment.AddInstruction("br label %v", AsLocalVariable(conditionExprLabel))

	fragment.AddLabel(forEndLabel)
}

func (visitor *CodegenVisitor) VisitEnterSwitchStmtNode(node *ast.SwitchStmtNode) {
	node.EndLabel = visitor.labeller.NewSet("switch", "end")
}

func (visitor *CodegenVisitor) VisitLeaveSwitchStmtNode(node *ast.SwitchStmtNode) {
	testTyping := node.TestExpr.GetTyping()
	cases := len(node.CaseExprs)

	fragment := visitor.newVoidCode(node)

	startLabel := visitor.labeller.NewSet("switch", "start")
	endLabel := node.EndLabel

	caseExprLabels := make([]string, cases)
	caseBlockLabels := make([]string, cases)
	defaultBlockLabel := visitor.labeller.Label("switch", "defaultBlock")

	for i := 0; i < cases; i++ {
		caseExprLabels[i] = visitor.labeller.Label("switch", "caseExpr", strconv.Itoa(i))
		caseBlockLabels[i] = visitor.labeller.Label("switch", "caseBlock", strconv.Itoa(i))
	}

	fragment.AddInstruction("br label %v", AsLocalVariable(startLabel))
	fragment.AddLabel(startLabel)

	testExprFragment := visitor.removeValueCode(node.TestExpr)
	testExprResult := testExprFragment.GetResult()

	fragment.Append(testExprFragment)

	// cases
	for i := 0; i < cases; i++ {
		caseExprLabel := caseExprLabels[i]
		caseBlockLabel := caseBlockLabels[i]

		fragment.AddInstruction("br label %v", AsLocalVariable(caseExprLabel))
		fragment.AddLabel(caseExprLabel)

		caseExprFragment := visitor.removeValueCode(node.CaseExprs[i])
		caseExprResult := caseExprFragment.GetResult()

		fragment.Append(caseExprFragment)

		operatorCodegen := NewOperatorCodegen(nil, signature.SHALLOW_EQUAL, testTyping, nil, nil)

		equalOperator := operatorCodegen.GenerateComparisonOpcode()

		comparisonResult := fragment.AddOperation("%v %v %v, %v", equalOperator, testTyping.IrType(), testExprResult, caseExprResult)
		truthResult := fragment.AddOperation("icmp eq i1 %v, 1", comparisonResult)

		// If the current case passes, it should jump to the next non-empty block
		// If the switch statement does not have a non-empty block, then it should jump to the endLabel

		var jumpLabelWhenTruthy string

		nextNonEmptyBlockIndex := node.FindTheNextNonEmptyBlockIndexAt(i)

		if nextNonEmptyBlockIndex < cases {
			jumpLabelWhenTruthy = caseBlockLabels[nextNonEmptyBlockIndex]
		} else if !node.IsEmptyDefaultBlock() {
			jumpLabelWhenTruthy = defaultBlockLabel
		} else {
			jumpLabelWhenTruthy = endLabel
		}

		var jumpLabelWhenFalsy string

		if i+1 < cases {
			jumpLabelWhenFalsy = caseExprLabels[i+1]
		} else if !node.IsEmptyDefaultBlock() {
			jumpLabelWhenFalsy = defaultBlockLabel
		} else {
			jumpLabelWhenFalsy = endLabel
		}

		fragment.AddInstruction("br i1 %v, label %v, label %v", truthResult, AsLocalVariable(jumpLabelWhenTruthy), AsLocalVariable(jumpLabelWhenFalsy))

		if !node.IsEmptyCaseBlockAt(i) {
			fragment.AddLabel(caseBlockLabel)
			fragment.Append(visitor.removeVoidCode(node.CaseBlocks[i]))

			var jumpLabelAfterCaseBlock string

			nextNonEmptyBlockIndex := node.FindTheNextNonEmptyBlockIndexAt(i + 1)

			if nextNonEmptyBlockIndex < cases {
				jumpLabelAfterCaseBlock = caseBlockLabels[nextNonEmptyBlockIndex]
			} else if !node.IsEmptyDefaultBlock() {
				jumpLabelAfterCaseBlock = defaultBlockLabel
			} else {
				jumpLabelAfterCaseBlock = endLabel
			}

			fragment.AddInstruction("br label %v", AsLocalVariable(jumpLabelAfterCaseBlock))
		}
	}

	// default
	if !node.IsEmptyDefaultBlock() {
		fragment.AddLabel(defaultBlockLabel)
		fragment.Append(visitor.removeVoidCode(node.DefaultBlock))
		fragment.AddInstruction("br label %v", AsLocalVariable(endLabel))
	}

	fragment.AddLabel(endLabel)
}

func (visitor *CodegenVisitor) VisitBreakNode(node *ast.BreakNode) {
	fragment := visitor.newVoidCode(node)

	breakLabel := node.FindBreakLabel()

	fragment.AddInstruction("br label %v", AsLocalVariable(breakLabel))
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
	typing := node.Lhs.GetTyping()

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

	if node.Val == 0 {
		fragment.result = "0.0"
	} else {
		fragment.result = strconv.FormatFloat(float64(node.Val), 'f', 6, 32)
	}
}

// VisitCharacterNode do something
func (visitor *CodegenVisitor) VisitCharacterNode(node *ast.CharacterNode) {

}

// VisitStringNode do something
func (visitor *CodegenVisitor) VisitStringNode(node *ast.StringNode) {
	fragment := visitor.newAddressCode(node)

	stringValue := node.EscapeVal()
	stringLength := node.EscapedStringLength()

	stringGlobalIdentifier := visitor.constants.AddOperation("private constant [%v x i8] c\"%v\", align 1", stringLength, stringValue)

	fragment.AddOperation("getelementptr inbounds [%v x i8], [%v x i8]* %v, i32 0, i32 0", stringLength, stringLength, stringGlobalIdentifier)
}

// VisitIdentifierNode do something
func (visitor *CodegenVisitor) VisitIdentifierNode(node *ast.IdentifierNode) {
	fragment := visitor.newAddressCode(node)

	identifier := node.LocalIdentifier()

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
