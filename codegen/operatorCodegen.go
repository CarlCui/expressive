package codegen

import (
	"fmt"

	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/typing"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

type OperatorCodegen struct {
	fragment   *BlocksFragment
	operands   []Fragment
	operator   signature.Operator
	typing     typing.Typing
	labeller   *Labeller
	compIPreds map[string]enum.IPred
	compFPreds map[string]enum.FPred
}

func NewOperatorCodegen(fragment *BlocksFragment, operator signature.Operator, typing typing.Typing, labeller *Labeller, operands ...Fragment) *OperatorCodegen {
	compIPreds := map[string]enum.IPred{
		"eq":  enum.IPredEQ,
		"ne":  enum.IPredNE,
		"sge": enum.IPredSGE,
		"sgt": enum.IPredSGT,
		"sle": enum.IPredSLE,
		"slt": enum.IPredSLT,
		"uge": enum.IPredUGE,
		"ugt": enum.IPredUGT,
		"ule": enum.IPredULE,
		"ult": enum.IPredULT,
	}

	// only implemented ordered ones
	compFPreds := map[string]enum.FPred{
		"false": enum.FPredFalse,
		"oeq":   enum.FPredOEQ,
		"oge":   enum.FPredOGE,
		"ogt":   enum.FPredOGT,
		"ole":   enum.FPredOLE,
		"olt":   enum.FPredOLT,
		"one":   enum.FPredONE,
		"ord":   enum.FPredORD,
	}

	return &OperatorCodegen{
		fragment,
		operands,
		operator,
		typing,
		labeller,
		compIPreds,
		compFPreds,
	}
}

func (gen *OperatorCodegen) GenerateCode() {
	switch gen.operator {
	case signature.ADD:
		gen.generateAdd()
	case signature.SUBTRACT:
		gen.generateSubtract()
	case signature.MULTIPLY:
		gen.generateMultiply()
	case signature.DIVIDE:
		gen.generateDivide()
	case signature.MODULO:
		gen.generateRemainder()
	case signature.LOGIC_AND:
		gen.generateLogicalAnd()
	case signature.LOGIC_OR:
		gen.generateLogicalOr()
	case signature.LOGIC_NOT:
		gen.generateLogicalNot()
	case signature.IF_ELSE:
		gen.generateIfElse()
	case signature.GREATER, signature.GREATER_OR_EQUAL,
		signature.LESS, signature.LESS_OR_EQUAL,
		signature.SHALLOW_EQUAL, signature.SHALLOW_NOT_EQUAL,
		signature.DEEP_EQUAL, signature.DEEP_NOT_EQUAL:
		gen.generateComparison()
	default:
		panic(fmt.Sprintf("Not implemented for operator %v", gen.operator))
	}
}

func (gen *OperatorCodegen) GenerateComparisonInstr(frag *BlocksFragment) func(value.Value, value.Value) value.Value {

	var conditionCodePrefix string // signed, unsigned, ordered or unordered depending on the typing

	var opcode string
	switch gen.typing {
	case typing.INT:
		conditionCodePrefix = "s" // signed
	case typing.BOOL:
		conditionCodePrefix = "u" // does not matter if signed or unsigned
	case typing.FLOAT:
		conditionCodePrefix = "o" // ordered, since neither can be QNAN
	// TODO: implement string and char
	default:
		gen.panicOnMismatchCodegen()
	}

	switch gen.operator {
	case signature.GREATER:
		opcode += conditionCodePrefix + "gt"
	case signature.GREATER_OR_EQUAL:
		opcode += conditionCodePrefix + "ge"
	case signature.LESS:
		opcode += conditionCodePrefix + "lt"
	case signature.LESS_OR_EQUAL:
		opcode += conditionCodePrefix + "le"
	case signature.SHALLOW_EQUAL:
		if gen.typing == typing.FLOAT {
			opcode += conditionCodePrefix + "eq"
		} else {
			opcode += "eq"
		}
	case signature.SHALLOW_NOT_EQUAL:
		if gen.typing == typing.FLOAT {
			opcode += conditionCodePrefix + "ne"
		} else {
			opcode += "ne"
		}
	case signature.DEEP_EQUAL:
		if gen.typing == typing.FLOAT {
			opcode += conditionCodePrefix + "eq"
		} else {
			opcode += "eq"
		}
	case signature.DEEP_NOT_EQUAL:
		if gen.typing == typing.FLOAT {
			opcode += conditionCodePrefix + "ne"
		} else {
			opcode += "ne"
		}
	default:
		gen.panicOnMismatchCodegen()
	}

	switch gen.typing {
	case typing.INT, typing.BOOL:
		predicate := gen.compIPreds[opcode]
		return func(op1, op2 value.Value) value.Value {
			return ir.NewICmp(predicate, op1, op2)
		}
	case typing.FLOAT:
		predicate := gen.compFPreds[opcode]
		return func(op1, op2 value.Value) value.Value {
			return ir.NewFCmp(predicate, op1, op2)
		}
	default:
		panic("does not support comparison on type %v " + gen.typing.String())
	}
}

func (gen *OperatorCodegen) generateComparison() {

	instr := gen.GenerateComparisonInstr(gen.fragment)

	gen.generateBinary(instr)
}

func (gen *OperatorCodegen) generateBinary(instrFunction func(value.Value, value.Value) value.Value) {
	frag := gen.fragment

	if len(frag.Blocks) == 0 {
		frag.NewBlock("")
	}

	gen.checkOperandsLength(2)

	frag1 := gen.operands[0]
	frag2 := gen.operands[1]

	op1 := frag1.GetResult()
	op2 := frag2.GetResult()

	frag.Append(frag1)
	frag.Append(frag2)

	result := instrFunction(op1, op2)
	instr := result.(ir.Instruction)
	frag.CurrentBlock.Insts = append(frag.CurrentBlock.Insts, instr)

	frag.resultValue = result
}

func (gen *OperatorCodegen) generateAdd() {
	var instr func(value.Value, value.Value) value.Value

	switch gen.typing {
	case typing.INT:
		instr = func(op1, op2 value.Value) value.Value {
			return ir.NewAdd(op1, op2)
		}
	case typing.FLOAT:
		instr = func(op1, op2 value.Value) value.Value {
			return ir.NewFAdd(op1, op2)
		}
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(instr)
}

func (gen *OperatorCodegen) generateSubtract() {
	var instr func(value.Value, value.Value) value.Value

	switch gen.typing {
	case typing.INT:
		instr = func(op1, op2 value.Value) value.Value {
			return ir.NewSub(op1, op2)
		}
	case typing.FLOAT:
		instr = func(op1, op2 value.Value) value.Value {
			return ir.NewFSub(op1, op2)
		}
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(instr)
}

func (gen *OperatorCodegen) generateMultiply() {
	var instr func(value.Value, value.Value) value.Value
	switch gen.typing {
	case typing.INT:
		instr = func(op1, op2 value.Value) value.Value {
			return ir.NewMul(op1, op2)
		}
	case typing.FLOAT:
		instr = func(op1, op2 value.Value) value.Value {
			return ir.NewFMul(op1, op2)
		}
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(instr)
}

func (gen *OperatorCodegen) generateDivide() {
	var instr func(value.Value, value.Value) value.Value
	switch gen.typing {
	case typing.INT:
		instr = func(op1, op2 value.Value) value.Value {
			return ir.NewSDiv(op1, op2)
		}
	case typing.FLOAT:
		instr = func(op1, op2 value.Value) value.Value {
			return ir.NewFDiv(op1, op2)
		}
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(instr)
}

func (gen *OperatorCodegen) generateRemainder() {
	var instr func(value.Value, value.Value) value.Value
	switch gen.typing {
	case typing.INT:
		instr = func(op1, op2 value.Value) value.Value {
			return ir.NewSRem(op1, op2)
		}
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(instr)
}

func (gen *OperatorCodegen) generateLogicalAnd() {

	switch gen.typing {
	case typing.BOOL:

	default:
		gen.panicOnMismatchCodegen()
	}

	gen.checkOperandsLength(2)

	frag := gen.fragment

	frag1 := gen.operands[0]
	frag2 := gen.operands[1]

	op1 := frag1.GetResult()
	op2 := frag2.GetResult()

	entry := gen.labeller.NewSet("LAND", "entry")
	condTrue := gen.labeller.Label("LAND", "cond", "true")
	condFalse := gen.labeller.Label("LAND", "cond", "false")
	condFalseEval := gen.labeller.Label("LAND", "cond", "false", "eval")

	condTrueBlock := ir.NewBlock(condTrue)
	condFalseBlock := ir.NewBlock(condFalse)
	condFalseEvalBlock := ir.NewBlock(condFalseEval)

	frag.Append(frag1)

	frag.NewBlock(entry)
	entryBlock := frag.CurrentBlock

	frag1IsTrue := frag.CurrentBlock.NewICmp(enum.IPredEQ, op1, constant.False)
	frag.CurrentBlock.NewCondBr(frag1IsTrue, condTrueBlock, condFalseBlock)

	frag.AddBlock(condFalseBlock)
	frag.Append(frag2)

	frag.AddBlock(condFalseEvalBlock)
	instrAnd := frag.CurrentBlock.NewAnd(op1, op2)

	frag.AddBlock(condTrueBlock)
	instrPhi := frag.CurrentBlock.NewPhi(ir.NewIncoming(constant.False, entryBlock), ir.NewIncoming(instrAnd, condFalseEvalBlock))

	frag.resultValue = instrPhi
}

func (gen *OperatorCodegen) generateLogicalOr() {

	switch gen.typing {
	case typing.BOOL:

	default:
		gen.panicOnMismatchCodegen()
	}

	gen.checkOperandsLength(2)

	frag1 := gen.operands[0]
	frag2 := gen.operands[1]

	op1 := frag1.GetResult()
	op2 := frag2.GetResult()

	entry := gen.labeller.NewSet("LOR", "entry")
	condTrue := gen.labeller.Label("LOR", "cond", "true")
	condFalse := gen.labeller.Label("LOR", "cond", "false")
	condFalseEval := gen.labeller.Label("LOR", "cond", "false", "eval")

	entryBlock := ir.NewBlock(entry)
	condTrueBlock := ir.NewBlock(condTrue)
	condFalseBlock := ir.NewBlock(condFalse)
	condFalseEvalBlock := ir.NewBlock(condFalseEval)

	frag := gen.fragment

	frag.Append(frag1)

	frag.AddBlock(entryBlock)

	frag1IsTrue := frag.CurrentBlock.NewICmp(enum.IPredEQ, op1, constant.True)
	frag.CurrentBlock.NewCondBr(frag1IsTrue, condTrueBlock, condFalseBlock)

	frag.AddBlock(condFalseBlock)
	frag.Append(frag2)

	frag.AddBlock(condFalseEvalBlock)

	result := frag.CurrentBlock.NewOr(op1, op2)

	frag.CurrentBlock.NewBr(condTrueBlock)

	frag.AddBlock(condTrueBlock)

	phiInstr := frag.CurrentBlock.NewPhi(ir.NewIncoming(constant.True, entryBlock), ir.NewIncoming(result, condFalseEvalBlock))

	frag.resultValue = phiInstr
}

func (gen *OperatorCodegen) generateLogicalNot() {

	switch gen.typing {
	case typing.BOOL:

	default:
		gen.panicOnMismatchCodegen()
	}

	gen.checkOperandsLength(1)

	frag1 := gen.operands[0]

	op1 := frag1.GetResult()

	frag := gen.fragment
	frag.NewBlock("")

	frag.Append(frag1)

	result := frag.CurrentBlock.NewXor(op1, constant.True)

	frag.resultValue = result
}

func (gen *OperatorCodegen) generateIfElse() {

	gen.checkOperandsLength(3)

	entry := ir.NewBlock(gen.labeller.NewSet("ifelse", "entry"))
	condTrue := ir.NewBlock(gen.labeller.Label("ifelse", "cond", "true"))
	condTrueEval := ir.NewBlock(gen.labeller.Label("ifelse", "cond", "true", "eval"))
	condFalse := ir.NewBlock(gen.labeller.Label("ifelse", "cond", "false"))
	condFalseEval := ir.NewBlock(gen.labeller.Label("ifelse", "cond", "false", "eval"))
	condEnd := ir.NewBlock(gen.labeller.Label("ifelse", "cond", "end"))

	frag1 := gen.operands[0]
	frag2 := gen.operands[1]
	frag3 := gen.operands[2]

	op1 := frag1.GetResult()
	op2 := frag2.GetResult()
	op3 := frag3.GetResult()

	frag := gen.fragment

	frag.AddBlock(entry)

	gen.fragment.Append(frag1)

	comp := frag.CurrentBlock.NewICmp(enum.IPredEQ, op1, constant.True)
	frag.CurrentBlock.NewCondBr(comp, condTrue, condFalse)

	frag.AddBlock(condTrue)
	frag.Append(frag2)

	frag.AddBlock(condTrueEval)
	frag.CurrentBlock.NewBr(condEnd)

	frag.AddBlock(condFalse)
	gen.fragment.Append(frag3)

	frag.AddBlock(condFalseEval)
	frag.CurrentBlock.NewBr(condEnd)

	frag.AddBlock(condEnd)
	phiInstr := frag.CurrentBlock.NewPhi(ir.NewIncoming(op2, condTrueEval), ir.NewIncoming(op3, condFalseEval))

	frag.resultValue = phiInstr
}

func (gen *OperatorCodegen) checkOperandsLength(needed int) {
	if len(gen.operands) != needed {
		gen.panicOnMismatchOperands(needed, len(gen.operands))
	}
}

func (gen *OperatorCodegen) panicOnMismatchOperands(needed int, actual int) {
	panic(fmt.Sprintf("Not enough operands: needed %v, actual %v", needed, actual))
}

func (gen *OperatorCodegen) panicOnMismatchCodegen() {
	panic(fmt.Sprintf("Code cannot be generated with operator %v and typing %v", gen.operator, gen.typing))
}
