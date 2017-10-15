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

	default:
		panic(fmt.Sprintf("Not implemented for operator %v", gen.operator))
	}
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
	condEnd := gen.labeller.Label("LAND", "cond", "end")

	gen.fragment.Append(frag1)
	frag1IsTrue := gen.fragment.AddOperation("icmp eq i1 %v, 0", localIdentifier1)
	gen.fragment.AddInstruction("br i1 %v, label %%v, label %%v", frag1IsTrue, condTrue, condFalse)

	gen.fragment.AddLabel(condFalse)
	gen.fragment.Append(frag2)
	result := gen.fragment.AddOperation("and i1 %v, %v", localIdentifier1, localIdentifier2)

	gen.fragment.AddLabel(condTrue)
	gen.fragment.AddOperation("phi i1 [false, %%v], [%v, %%v]", entry, result, condFalse)
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
	condEnd := gen.labeller.Label("LOR", "cond", "end")

	gen.fragment.Append(frag1)
	frag1IsTrue := gen.fragment.AddOperation("icmp eq i1 %v, 1", localIdentifier1)
	gen.fragment.AddInstruction("br i1 %v, label %%v, label %%v", frag1IsTrue, condTrue, condFalse)

	gen.fragment.AddLabel(condFalse)
	gen.fragment.Append(frag2)
	result := gen.fragment.AddOperation("or i1 %v, %v", localIdentifier1, localIdentifier2)

	gen.fragment.AddLabel(condTrue)
	gen.fragment.AddOperation("phi i1 [true, %%v], [%v, %%v]", entry, result, condFalse)
	gen.fragment.AddLabel(condEnd)
}

func (gen *OperatorCodegen) generateIfElse() {

	gen.checkOperandsLength(3)

	entry := gen.labeller.NewSet("ifelse", "entry")
	condTrue := gen.labeller.Label("ifelse", "cond", "true")
	condFalse := gen.labeller.Label("ifelse", "cond", "false")
	condEnd := gen.labeller.Label("ifelse", "cond", "end")

	frag1 := gen.operands[0]
	frag2 := gen.operands[1]
	frag3 := gen.operands[2]

	localIdentifier1 := frag1.GetResult()
	localIdentifier2 := frag2.GetResult()
	localIdentifier3 := frag3.GetResult()

	gen.fragment.AddLabel(entry)
	gen.fragment.Append(frag1)
	comp := gen.fragment.AddOperation("icmp eq i1 %v, 1", localIdentifier1)
	gen.fragment.AddInstruction("br i1 %v, label %%v, label %%v", comp, condTrue, condFalse)

	gen.fragment.AddLabel(condTrue)
	gen.fragment.Append(frag2)
	gen.fragment.AddInstruction("br label %%v", condEnd)

	gen.fragment.AddLabel(condFalse)
	gen.fragment.Append(frag3)
	gen.fragment.AddInstruction("br label %v", condEnd)

	gen.fragment.AddLabel(condEnd)
	gen.fragment.AddOperation("phi %v [%v, %%v], [%v, %%v]", gen.typing.IrType(), localIdentifier2, condTrue, localIdentifier3, condFalse)
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
