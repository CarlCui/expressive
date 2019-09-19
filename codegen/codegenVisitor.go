package codegen

import (
	"fmt"
	"strconv"

	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/typing"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/logger"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// CodegenVisitor visits each node and generates llvm IR.
type CodegenVisitor struct {
	logger                  logger.Logger
	labeller                *Labeller
	constants               []*ir.Global // global constants
	externals               []*ir.Func   // external function declarations
	codeMap                 map[ast.Node]Fragment
	globalIdentifierTracker *GlobalIdentifierTracker
}

// Init with a logger
func (visitor *CodegenVisitor) Init(logger logger.Logger) {
	visitor.logger = logger
	visitor.labeller = &Labeller{0}
	visitor.constants = make([]*ir.Global, 0)
	visitor.externals = make([]*ir.Func, 0)
	visitor.codeMap = make(map[ast.Node]Fragment)
	visitor.globalIdentifierTracker = &GlobalIdentifierTracker{index: 0}
}

func (visitor *CodegenVisitor) checkIfFragmentExists(node ast.Node) {
	if _, exists := visitor.codeMap[node]; exists {
		panic(fmt.Sprintf("Code for node %v already exists.", node))
	}
}

func (visitor *CodegenVisitor) newModuleFragment(node ast.Node) *ModuleFragment {
	visitor.checkIfFragmentExists(node)
	frag := NewModuleFragment(nil)
	visitor.codeMap[node] = frag

	return frag
}

func (visitor *CodegenVisitor) newFunctionsFragment(node ast.Node) *FunctionsFragment {
	visitor.checkIfFragmentExists(node)
	frag := NewFunctionsFragment()
	visitor.codeMap[node] = frag

	return frag
}

func (visitor *CodegenVisitor) newBlocksFragment(node ast.Node, resultType ResultType) *BlocksFragment {
	visitor.checkIfFragmentExists(node)
	frag := NewBlocksFragment(resultType)
	visitor.codeMap[node] = frag

	return frag
}

func (visitor *CodegenVisitor) getAndRemoveFragment(node ast.Node) Fragment {
	fragment, exists := visitor.codeMap[node]

	if !exists {
		panic(fmt.Sprintf("Code for node %v does not exist.", node))
	}

	delete(visitor.codeMap, node)

	return fragment
}

func (visitor *CodegenVisitor) removeVoidFragment(node ast.Node) Fragment {
	fragment := visitor.getAndRemoveFragment(node)

	if fragment.GetResultType() != VOID {
		panic(fmt.Sprintf("Code fragment does not produce void result: %v", node))
	}

	return fragment
}

func (visitor *CodegenVisitor) removePointerFragment(node ast.Node) Fragment {
	fragment := visitor.getAndRemoveFragment(node)

	if fragment.GetResultType() != POINTER {
		panic(fmt.Sprintf("Code fragment does not produce pointer result: %v", node))
	}

	return fragment
}

func (visitor *CodegenVisitor) removeValueFragment(node ast.Node) Fragment {
	fragment := visitor.getAndRemoveFragment(node)

	if fragment.GetResultType() != VALUE && fragment.GetResultType() != POINTER {
		panic(fmt.Sprintf("Code fragment does not produce value result: %v", node))
	}

	if fragment.GetResultType() == POINTER {
		visitor.dereferencePointer(node, fragment)
	}

	return fragment
}

func (visitor *CodegenVisitor) dereferencePointer(node ast.Node, fragment Fragment) {
	switch f := fragment.(type) {
	case *ModuleFragment:
	case *FunctionsFragment:
		panic("Fragment should be of void result type")
	case *BlocksFragment:
		f.resultType = VALUE

		result := f.GetResult()

		if len(f.Blocks) == 0 {
			f.NewBlock("")
		}

		load := f.CurrentBlock.NewLoad(result)
		f.resultValue = load
	}
}

// VisitEnterProgramNode creates program scope
func (visitor *CodegenVisitor) VisitEnterProgramNode(node *ast.ProgramNode) {

	printfDeclaration := ir.NewFunc("printf", types.I32, ir.NewParam("", types.I8Ptr))
	printfDeclaration.Sig.Variadic = true
	printfDeclaration.FuncAttrs = append(printfDeclaration.FuncAttrs, enum.FuncAttrNoUnwind)

	visitor.externals = append(visitor.externals, printfDeclaration)
}

// VisitLeaveProgramNode closes program scope
func (visitor *CodegenVisitor) VisitLeaveProgramNode(node *ast.ProgramNode) {
	// generates main function
	fragment := visitor.newModuleFragment(node)

	for _, external := range visitor.externals {
		external.Parent = fragment.Module
	}

	fragment.Module.Funcs = append(fragment.Module.Funcs, visitor.externals...)

	mainFunc := ir.NewFunc("main", types.I32)

	mainFunc.Blocks = make([]*ir.Block, 0)

	functionsFragment := NewFunctionsFragment()
	functionsFragment.AddFunc(mainFunc)

	for _, child := range node.Chilren {
		functionsFragment.Append(visitor.removeVoidFragment(child))
	}

	var lastBlock *ir.Block

	numberOfBlocks := len(mainFunc.Blocks)

	if numberOfBlocks == 0 {
		lastBlock = mainFunc.NewBlock("")
	} else {
		lastBlock = mainFunc.Blocks[numberOfBlocks-1]
	}

	zeroConstant := constant.NewInt(types.I32, 0)

	lastBlock.NewRet(zeroConstant)

	fragment.Append(functionsFragment)
}

func (visitor *CodegenVisitor) VisitEnterBlockNode(node *ast.BlockNode) {

}

func (visitor *CodegenVisitor) VisitLeaveBlockNode(node *ast.BlockNode) {
	fragment := visitor.newBlocksFragment(node, VOID)

	for _, child := range node.Stmts {
		fragment.Append(visitor.removeVoidFragment(child))
	}
}

// stmts

// VisitEnterVariableDeclarationNode do something
func (visitor *CodegenVisitor) VisitEnterVariableDeclarationNode(node *ast.VariableDeclarationNode) {

}

// VisitLeaveVariableDeclarationNode do something
func (visitor *CodegenVisitor) VisitLeaveVariableDeclarationNode(node *ast.VariableDeclarationNode) {
	fragment := visitor.newBlocksFragment(node, VOID)

	identifierNode := node.Identifier.(*ast.IdentifierNode)
	identifierTyping := identifierNode.GetTyping()
	irType := identifierTyping.IrType()

	fragment.NewBlock("")
	allocaInstr := fragment.CurrentBlock.NewAlloca(irType)
	allocaInstr.SetName(identifierNode.LocalIdentifier())

	if node.Expr == nil {
		// load default value
		var defaultValue value.Value

		// TODO: finish default values for types
		switch identifierTyping {
		case typing.INT:
			defaultValue = constant.NewInt(typing.INT.IrType().(*types.IntType), 0)
		case typing.FLOAT:
			defaultValue = constant.NewFloat(typing.FLOAT.IrType().(*types.FloatType), 0)
		}

		fragment.CurrentBlock.NewStore(defaultValue, allocaInstr)
	} else {
		exprFragment := visitor.removeValueFragment(node.Expr)

		exprResult := exprFragment.GetResult()

		fragment.Append(exprFragment)

		fragment.CurrentBlock.NewStore(exprResult, allocaInstr)
	}
}

// VisitEnterAssignmentNode do something
func (visitor *CodegenVisitor) VisitEnterAssignmentNode(node *ast.AssignmentNode) {
}

// VisitLeaveAssignmentNode do something
func (visitor *CodegenVisitor) VisitLeaveAssignmentNode(node *ast.AssignmentNode) {
	fragment := visitor.newBlocksFragment(node, VOID)
	fragment.NewBlock("")

	identifierNode := node.Identifier.(*ast.IdentifierNode)
	irType := identifierNode.GetTyping().IrType()
	variableName := identifierNode.LocalIdentifier()

	exprFragment := visitor.removeValueFragment(node.Expr)

	exprResult := exprFragment.GetResult()

	fragment.Append(exprFragment)

	allocaInstr := ir.NewAlloca(irType)
	allocaInstr.SetName(variableName)

	fragment.CurrentBlock.NewStore(exprResult, allocaInstr)
}

// VisitEnterPrintNode do something
func (visitor *CodegenVisitor) VisitEnterPrintNode(node *ast.PrintNode) {

}

// VisitLeavePrintNode do something
func (visitor *CodegenVisitor) VisitLeavePrintNode(node *ast.PrintNode) {
	fragment := visitor.newBlocksFragment(node, VOID)
	fragment.NewBlock("")

	stringExprFrag := visitor.removePointerFragment(node.StringExpr)

	instructionArgs := make([]interface{}, 0)

	stringExprResult := stringExprFrag.GetResult()

	instructionArgs = append(instructionArgs, stringExprResult)

	fragment.Append(stringExprFrag)

	argResults := make([]value.Value, 0)

	for _, arg := range node.Args {
		argFrag := visitor.removeValueFragment(arg)
		argResult := argFrag.GetResult()

		fragment.Append(argFrag)

		argResults = append(argResults, argResult)
	}

	args := append([]value.Value{stringExprResult}, argResults...)

	fragment.CurrentBlock.NewCall(visitor.externals[0], args...)
}

func (visitor *CodegenVisitor) VisitEnterIfStmtNode(node *ast.IfStmtNode) {

}

func (visitor *CodegenVisitor) VisitLeaveIfStmtNode(node *ast.IfStmtNode) {
	fragment := visitor.newBlocksFragment(node, VOID)

	fragment.NewBlock(visitor.labeller.NewSet("if", "start"))
	ifEnd := ir.NewBlock(visitor.labeller.NewSet("if", "end"))
	ifElse := ir.NewBlock(visitor.labeller.Label("if", "else"))

	numberOfConditions := len(node.ConditionExprs)

	ifConditions := make([]*ir.Block, numberOfConditions)
	ifBlocks := make([]*ir.Block, numberOfConditions)

	for i := range ifConditions {
		ifConditions[i] = ir.NewBlock(visitor.labeller.Label("if", "condition", strconv.Itoa(i)))
		ifBlocks[i] = ir.NewBlock(visitor.labeller.Label("if", "block", strconv.Itoa(i)))
	}

	for i, conditionExpr := range node.ConditionExprs {
		condition := ifConditions[i]
		block := ifBlocks[i]

		exprFragment := visitor.removeValueFragment(conditionExpr)

		fragment.AddBlock(condition)

		fragment.Append(exprFragment)

		exprResult := exprFragment.GetResult()

		conditionResult := fragment.CurrentBlock.NewICmp(enum.IPredEQ, exprResult, constant.True)

		finalBlock := ifEnd

		if i+1 != numberOfConditions {
			finalBlock = ifConditions[i+1]
		} else if node.ElseBlock != nil {
			finalBlock = ifElse
		}

		fragment.CurrentBlock.NewCondBr(conditionResult, block, finalBlock)

		fragment.AddBlock(block)

		fragment.Append(visitor.removeVoidFragment(node.ConditionBlocks[i]))

		fragment.NewBlock("")
		fragment.CurrentBlock.NewBr(ifEnd)
	}

	if node.ElseBlock != nil {
		fragment.AddBlock(ifElse)
		fragment.Append(visitor.removeVoidFragment(node.ElseBlock))
	}

	fragment.AddBlock(ifEnd)
}

func (visitor *CodegenVisitor) VisitEnterWhileStmtNode(node *ast.WhileStmtNode) {
	node.EndBlock = ir.NewBlock(visitor.labeller.NewSet("while", "end"))
}

func (visitor *CodegenVisitor) VisitLeaveWhileStmtNode(node *ast.WhileStmtNode) {
	fragment := visitor.newBlocksFragment(node, VOID)

	whileStart := ir.NewBlock(visitor.labeller.NewSet("while", "start"))
	whileEnd := node.EndBlock
	whileBlock := ir.NewBlock(visitor.labeller.Label("while", "block"))

	fragment.AddBlock(whileStart)

	conditionExprFragment := visitor.removeValueFragment(node.ConditionExpr)

	fragment.Append(conditionExprFragment)

	conditionResult := fragment.CurrentBlock.NewICmp(enum.IPredEQ, conditionExprFragment.GetResult(), constant.True)

	fragment.CurrentBlock.NewCondBr(conditionResult, whileBlock, whileEnd)

	fragment.AddBlock(whileBlock)

	fragment.Append(visitor.removeVoidFragment(node.Block))

	fragment.NewBlock("")
	fragment.CurrentBlock.NewBr(whileStart)

	fragment.AddBlock(whileEnd)
}

func (visitor *CodegenVisitor) VisitEnterForStmtNode(node *ast.ForStmtNode) {
	node.EndBlock = ir.NewBlock(visitor.labeller.NewSet("for", "end"))
}

func (visitor *CodegenVisitor) VisitEnterForStmtNodeBeforeBlockNode(node *ast.ForStmtNode) {

}

func (visitor *CodegenVisitor) VisitLeaveForStmtNode(node *ast.ForStmtNode) {
	fragment := visitor.newBlocksFragment(node, VOID)

	forStart := ir.NewBlock(visitor.labeller.NewSet("for", "start"))
	forEnd := node.EndBlock

	conditionExpr := ir.NewBlock(visitor.labeller.Label("for", "condition"))
	iterationStmt := ir.NewBlock(visitor.labeller.Label("for", "iteration"))
	block := ir.NewBlock(visitor.labeller.Label("for", "block"))

	fragment.AddBlock(forStart)

	if node.InitializationStmt != nil {
		fragment.Append(visitor.removeVoidFragment(node.InitializationStmt))
	}

	fragment.AddBlock(conditionExpr)

	if node.ConditionExpr != nil {
		conditionExprFragment := visitor.removeValueFragment(node.ConditionExpr)

		fragment.Append(conditionExprFragment)

		conditionResult := fragment.CurrentBlock.NewICmp(enum.IPredEQ, conditionExprFragment.GetResult(), constant.True)

		fragment.CurrentBlock.NewCondBr(conditionResult, block, forEnd)
	}

	fragment.AddBlock(block)

	fragment.Append(visitor.removeVoidFragment(node.Block))

	fragment.AddBlock(iterationStmt)

	if node.IterationStmt != nil {
		fragment.Append(visitor.removeVoidFragment(node.IterationStmt))
	}

	fragment.NewBlock("")
	fragment.CurrentBlock.NewBr(conditionExpr)

	fragment.AddBlock(forEnd)
}

func (visitor *CodegenVisitor) VisitEnterSwitchStmtNode(node *ast.SwitchStmtNode) {
	node.EndBlock = ir.NewBlock(visitor.labeller.NewSet("switch", "end"))
}

func (visitor *CodegenVisitor) VisitLeaveSwitchStmtNode(node *ast.SwitchStmtNode) {
	testTyping := node.TestExpr.GetTyping()
	cases := len(node.CaseExprs)

	fragment := visitor.newBlocksFragment(node, VOID)

	start := ir.NewBlock(visitor.labeller.NewSet("switch", "start"))
	end := node.EndBlock

	caseExprs := make([]*ir.Block, cases)
	caseBlocks := make([]*ir.Block, cases)
	defaultBlock := ir.NewBlock(visitor.labeller.Label("switch", "defaultBlock"))

	for i := 0; i < cases; i++ {
		caseExprs[i] = ir.NewBlock(visitor.labeller.Label("switch", "caseExpr", strconv.Itoa(i)))
		caseBlocks[i] = ir.NewBlock(visitor.labeller.Label("switch", "caseBlock", strconv.Itoa(i)))
	}

	fragment.AddBlock(start)

	testExprFragment := visitor.removeValueFragment(node.TestExpr)
	testExprResult := testExprFragment.GetResult()

	fragment.Append(testExprFragment)

	// cases
	for i := 0; i < cases; i++ {
		caseExpr := caseExprs[i]
		caseBlock := caseBlocks[i]

		fragment.AddBlock(caseExpr)

		caseExprFragment := visitor.removeValueFragment(node.CaseExprs[i])
		caseExprResult := caseExprFragment.GetResult()

		fragment.Append(caseExprFragment)

		operatorCodegen := NewOperatorCodegen(nil, signature.SHALLOW_EQUAL, testTyping, nil, nil)

		equalOperator := operatorCodegen.GenerateComparisonInstr(fragment)

		comparisonResult := equalOperator(testExprResult, caseExprResult)

		instr := comparisonResult.(ir.Instruction)
		fragment.CurrentBlock.Insts = append(fragment.CurrentBlock.Insts, instr)

		truthResult := fragment.CurrentBlock.NewICmp(enum.IPredEQ, comparisonResult, constant.True)

		// If the current case passes, it should jump to the next non-empty block
		// If the switch statement does not have a non-empty block, then it should jump to the endLabel

		var jumpBlockWhenTruthy *ir.Block

		nextNonEmptyBlockIndex := node.FindTheNextNonEmptyBlockIndexAt(i)

		if nextNonEmptyBlockIndex < cases {
			jumpBlockWhenTruthy = caseBlocks[nextNonEmptyBlockIndex]
		} else if !node.IsEmptyDefaultBlock() {
			jumpBlockWhenTruthy = defaultBlock
		} else {
			jumpBlockWhenTruthy = end
		}

		var jumpBlockWhenFalsy *ir.Block

		if i+1 < cases {
			jumpBlockWhenFalsy = caseExprs[i+1]
		} else if !node.IsEmptyDefaultBlock() {
			jumpBlockWhenFalsy = defaultBlock
		} else {
			jumpBlockWhenFalsy = end
		}

		fragment.CurrentBlock.NewCondBr(truthResult, jumpBlockWhenTruthy, jumpBlockWhenFalsy)

		if !node.IsEmptyCaseBlockAt(i) {
			fragment.AddBlock(caseBlock)
			fragment.Append(visitor.removeVoidFragment(node.CaseBlocks[i]))

			var jumpBlockAfterCaseBlock *ir.Block

			nextNonEmptyBlockIndex := node.FindTheNextNonEmptyBlockIndexAt(i + 1)

			if nextNonEmptyBlockIndex < cases {
				jumpBlockAfterCaseBlock = caseBlocks[nextNonEmptyBlockIndex]
			} else if !node.IsEmptyDefaultBlock() {
				jumpBlockAfterCaseBlock = defaultBlock
			} else {
				jumpBlockAfterCaseBlock = end
			}

			fragment.NewBlock("")
			fragment.CurrentBlock.NewBr(jumpBlockAfterCaseBlock)
		}
	}

	// default
	if !node.IsEmptyDefaultBlock() {
		fragment.AddBlock(defaultBlock)
		fragment.Append(visitor.removeVoidFragment(node.DefaultBlock))
		fragment.NewBlock("")
		fragment.CurrentBlock.NewBr(end)
	}

	fragment.AddBlock(end)
}

func (visitor *CodegenVisitor) VisitBreakNode(node *ast.BreakNode) {
	fragment := visitor.newBlocksFragment(node, VOID)

	breakBlock := node.FindBreakBlock()

	fragment.NewBlock("")
	fragment.CurrentBlock.NewBr(breakBlock)
}

// exprs

// VisitEnterTernaryOperatorNode do something
func (visitor *CodegenVisitor) VisitEnterTernaryOperatorNode(node *ast.TernaryOperatorNode) {

}

// VisitLeaveTernaryOperatorNode do something
func (visitor *CodegenVisitor) VisitLeaveTernaryOperatorNode(node *ast.TernaryOperatorNode) {
	fragment := visitor.newBlocksFragment(node, VALUE)
	fragment.NewBlock("")

	fragment1 := visitor.removeValueFragment(node.Expr1)
	fragment2 := visitor.removeValueFragment(node.Expr2)
	fragment3 := visitor.removeValueFragment(node.Expr3)

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
	fragment := visitor.newBlocksFragment(node, VALUE)
	fragment.NewBlock("")

	fragment1 := visitor.removeValueFragment(node.Lhs)
	fragment2 := visitor.removeValueFragment(node.Rhs)

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
	fragment := visitor.newBlocksFragment(node, VALUE)
	fragment.NewBlock("")

	fragment1 := visitor.removeValueFragment(node.Expr)

	operator := node.Operator
	typing := node.GetTyping()

	operatorCodegen := NewOperatorCodegen(fragment, operator, typing, visitor.labeller, fragment1)

	operatorCodegen.GenerateCode()
}

// literal nodes

// VisitIntegerNode do something
func (visitor *CodegenVisitor) VisitIntegerNode(node *ast.IntegerNode) {
	fragment := visitor.newBlocksFragment(node, VALUE)

	irType := node.GetTyping().IrType()

	fragment.resultValue = constant.NewInt(irType.(*types.IntType), (int64)(node.Val))
}

// VisitFloatNode do something
func (visitor *CodegenVisitor) VisitFloatNode(node *ast.FloatNode) {
	fragment := visitor.newBlocksFragment(node, VALUE)

	irType := node.GetTyping().IrType()

	fragment.resultValue = constant.NewFloat(irType.(*types.FloatType), (float64)(node.Val))
}

// VisitCharacterNode do something
func (visitor *CodegenVisitor) VisitCharacterNode(node *ast.CharacterNode) {

}

// VisitStringNode do something
func (visitor *CodegenVisitor) VisitStringNode(node *ast.StringNode) {
	fragment := visitor.newBlocksFragment(node, POINTER)

	stringValue := node.StringValue()

	stringConstant := constant.NewCharArrayFromString(stringValue)
	stringGlobal := ir.NewGlobal(visitor.globalIdentifierTracker.NewIdentifier(), stringConstant.Type())
	stringGlobal.Init = stringConstant
	stringGlobal.Linkage = enum.LinkagePrivate
	stringGlobal.Align = 1

	visitor.constants = append(visitor.constants, stringGlobal)

	fragment.NewBlock("")
	result := fragment.CurrentBlock.NewGetElementPtr(stringGlobal, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	fragment.resultValue = result
}

// VisitIdentifierNode do something
func (visitor *CodegenVisitor) VisitIdentifierNode(node *ast.IdentifierNode) {
	fragment := visitor.newBlocksFragment(node, POINTER)

	identifier := node.LocalIdentifier()

	allocaInstr := ir.NewAlloca(node.GetTyping().IrType())
	allocaInstr.SetName(identifier)

	fragment.resultValue = allocaInstr

}

// VisitBooleanNode do something
func (visitor *CodegenVisitor) VisitBooleanNode(node *ast.BooleanNode) {
	fragment := visitor.newBlocksFragment(node, VALUE)

	var booleanValue value.Value

	if node.Val == true {
		booleanValue = constant.True
	} else {
		booleanValue = constant.False
	}

	fragment.resultValue = booleanValue
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
