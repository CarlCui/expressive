package codegen

// Instruction is one line of instruction in ir
type Instruction struct {
	Result    string // resulting local variable, if exists
	Operation string // actual operation
}

func (instruction *Instruction) Line() string {
	line := ""

	if len(instruction.Result) > 0 {
		line = line + instruction.Result + " = "
	}

	return line + instruction.Operation
}
