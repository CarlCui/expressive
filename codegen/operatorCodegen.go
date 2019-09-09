package codegen

import "github.com/carlcui/expressive/signature"
import "github.com/carlcui/expressive/typing"
import "fmt"

type OperatorCodegen struct {
	fragment *Fragment
	operands []*Fragment
	operator signature.Operator
	typing   typing.Typing
	labeller *Labeller
}

func NewOperatorCodegen(fragment *Fragment, operator signature.Operator, typing typing.Typing, labeller *Labeller, operands ...*Fragment) *OperatorCodegen {
	return &OperatorCodegen{
		fragment,
		operands,
		operator,
		typing,
		labeller,
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

func (gen *OperatorCodegen) generateComparison() {
	var opcode string

	var conditionCodePrefix string // signed, unsigned, ordered or unordered depending on the typing

	switch gen.typing {
	case typing.INT:
		opcode = "icmp "
		conditionCodePrefix = "s" // signed
	case typing.BOOL:
		opcode = "icmp "
		conditionCodePrefix = "u" // does not matter if signed or unsigned
	case typing.FLOAT:
		opcode = "fcmp "
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

	gen.generateBinary(opcode)
}

func (gen *OperatorCodegen) generateBinary(opcode string) {
	gen.checkOperandsLength(2)

	frag1 := gen.operands[0]
	frag2 := gen.operands[1]

	localIdentifier1 := frag1.GetResult()
	localIdentifier2 := frag2.GetResult()

	gen.fragment.Append(frag1)
	gen.fragment.Append(frag2)

	gen.fragment.AddOperation("%v %v %v, %v", opcode, gen.typing.IrType(), localIdentifier1, localIdentifier2)
}

func (gen *OperatorCodegen) generateAdd() {
	var opcode string
	switch gen.typing {
	case typing.INT:
		opcode = "add"
	case typing.FLOAT:
		opcode = "fadd"
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(opcode)
}

func (gen *OperatorCodegen) generateSubtract() {

	var opcode string
	switch gen.typing {
	case typing.INT:
		opcode = "sub"
	case typing.FLOAT:
		opcode = "fsub"
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(opcode)
}

func (gen *OperatorCodegen) generateMultiply() {

	var opcode string
	switch gen.typing {
	case typing.INT:
		opcode = "mul"
	case typing.FLOAT:
		opcode = "fmul"
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(opcode)
}

func (gen *OperatorCodegen) generateDivide() {

	var opcode string
	switch gen.typing {
	case typing.INT:
		opcode = "sdiv"
	case typing.FLOAT:
		opcode = "fdiv"
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(opcode)
}

func (gen *OperatorCodegen) generateRemainder() {

	var opcode string
	switch gen.typing {
	case typing.INT:
		opcode = "srem"
	default:
		gen.panicOnMismatchCodegen()
	}

	gen.generateBinary(opcode)
}

func (gen *OperatorCodegen) generateLogicalAnd() {

	switch gen.typing {
	case typing.BOOL:

	default:
		gen.panicOnMismatchCodegen()
	}

	gen.checkOperandsLength(2)

	frag1 := gen.operands[0]
	frag2 := gen.operands[1]

	localIdentifier1 := frag1.GetResult()
	localIdentifier2 := frag2.GetResult()

	entry := gen.labeller.NewSet("LAND", "entry")
	condTrue := gen.labeller.Label("LAND", "cond", "true")
	condFalse := gen.labeller.Label("LAND", "cond", "false")
	condFalseEval := gen.labeller.Label("LAND", "cond", "false", "eval")
	condEnd := gen.labeller.Label("LAND", "cond", "end")

	gen.fragment.Append(frag1)

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(entry))
	gen.fragment.AddLabel(entry)

	frag1IsTrue := gen.fragment.AddOperation("icmp eq i1 %v, 0", localIdentifier1)
	gen.fragment.AddInstruction("br i1 %v, label %v, label %v", frag1IsTrue, AsLocalVariable(condTrue), AsLocalVariable(condFalse))

	gen.fragment.AddLabel(condFalse)
	gen.fragment.Append(frag2)

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condFalseEval))
	gen.fragment.AddLabel(condFalseEval)

	result := gen.fragment.AddOperation("and i1 %v, %v", localIdentifier1, localIdentifier2)
	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condTrue))

	gen.fragment.AddLabel(condTrue)
	gen.fragment.AddOperation("phi i1 [false, %v], [%v, %v]", AsLocalVariable(entry), result, AsLocalVariable(condFalseEval))

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condEnd))
	gen.fragment.AddLabel(condEnd)
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

	localIdentifier1 := frag1.GetResult()
	localIdentifier2 := frag2.GetResult()

	entry := gen.labeller.NewSet("LOR", "entry")
	condTrue := gen.labeller.Label("LOR", "cond", "true")
	condFalse := gen.labeller.Label("LOR", "cond", "false")
	condFalseEval := gen.labeller.Label("LOR", "cond", "false", "eval")
	condEnd := gen.labeller.Label("LOR", "cond", "end")

	gen.fragment.Append(frag1)

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(entry))
	gen.fragment.AddLabel(entry)

	frag1IsTrue := gen.fragment.AddOperation("icmp eq i1 %v, 1", localIdentifier1)
	gen.fragment.AddInstruction("br i1 %v, label %v, label %v", frag1IsTrue, AsLocalVariable(condTrue), AsLocalVariable(condFalse))

	gen.fragment.AddLabel(condFalse)
	gen.fragment.Append(frag2)

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condFalseEval))
	gen.fragment.AddLabel(condFalseEval)

	result := gen.fragment.AddOperation("or i1 %v, %v", localIdentifier1, localIdentifier2)
	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condTrue))

	gen.fragment.AddLabel(condTrue)
	gen.fragment.AddOperation("phi i1 [true, %v], [%v, %v]", AsLocalVariable(entry), result, AsLocalVariable(condFalseEval))

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condEnd))
	gen.fragment.AddLabel(condEnd)
}

func (gen *OperatorCodegen) generateLogicalNot() {

	switch gen.typing {
	case typing.BOOL:

	default:
		gen.panicOnMismatchCodegen()
	}

	gen.checkOperandsLength(1)

	frag := gen.operands[0]

	fragResult := frag.GetResult()

	gen.fragment.Append(frag)

	gen.fragment.AddOperation("xor i1 %v, 1", fragResult)
}

func (gen *OperatorCodegen) generateIfElse() {

	gen.checkOperandsLength(3)

	entry := gen.labeller.NewSet("ifelse", "entry")
	condTrue := gen.labeller.Label("ifelse", "cond", "true")
	condTrueEval := gen.labeller.Label("ifelse", "cond", "true", "eval")
	condFalse := gen.labeller.Label("ifelse", "cond", "false")
	condFalseEval := gen.labeller.Label("ifelse", "cond", "false", "eval")
	condEnd := gen.labeller.Label("ifelse", "cond", "end")

	frag1 := gen.operands[0]
	frag2 := gen.operands[1]
	frag3 := gen.operands[2]

	localIdentifier1 := frag1.GetResult()
	localIdentifier2 := frag2.GetResult()
	localIdentifier3 := frag3.GetResult()

	gen.fragment.Append(frag1)

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(entry))
	gen.fragment.AddLabel(entry)

	comp := gen.fragment.AddOperation("icmp eq i1 %v, 1", localIdentifier1)
	gen.fragment.AddInstruction("br i1 %v, label %v, label %v", comp, AsLocalVariable(condTrue), AsLocalVariable(condFalse))

	gen.fragment.AddLabel(condTrue)
	gen.fragment.Append(frag2)

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condTrueEval))
	gen.fragment.AddLabel(condTrueEval)

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condEnd))

	gen.fragment.AddLabel(condFalse)
	gen.fragment.Append(frag3)

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condFalseEval))
	gen.fragment.AddLabel(condFalseEval)

	gen.fragment.AddInstruction("br label %v", AsLocalVariable(condEnd))

	gen.fragment.AddLabel(condEnd)
	gen.fragment.AddOperation("phi %v [%v, %v], [%v, %v]", gen.typing.IrType(), localIdentifier2, AsLocalVariable(condTrueEval), localIdentifier3, AsLocalVariable(condFalseEval))
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
