package codegen

import (
	"fmt"
)

type ResultType int

const (
	VOID ResultType = iota
	VALUE
	ADDRESS
)

// Fragment represents a fragment of instructions in asm
type Fragment struct {
	Lines             []IrLiner
	ResultType        ResultType
	identifierTracker IdentifierTracker
	result            string
}

// NewFragment is the initializer of Fragment
func NewFragment(resultType ResultType, identifierTracker IdentifierTracker) *Fragment {
	var fragment Fragment

	fragment.Lines = make([]IrLiner, 0)
	fragment.identifierTracker = identifierTracker
	fragment.ResultType = resultType
	fragment.result = ""

	return &fragment
}

// Append another code fragment
func (fragment *Fragment) Append(frag *Fragment) {
	fragment.Lines = append(fragment.Lines, frag.Lines...)
}

func (fragment *Fragment) AppendBefore(frag *Fragment) {
	fragment.Lines = append(frag.Lines, fragment.Lines...)
}

func (fragment *Fragment) AppendLines(lines ...IrLiner) {
	fragment.Lines = append(fragment.Lines, lines...)
}

// AddInstruction adds an instruction
func (fragment *Fragment) AddInstruction(format string, arguments ...interface{}) {
	operation := fmt.Sprintf(format, arguments...)

	instruction := &Instruction{Result: "", Operation: operation}

	fragment.Lines = append(fragment.Lines, instruction)
	fragment.result = instruction.Result
}

// AddOperation adds an operation, and return the result
func (fragment *Fragment) AddOperation(format string, arguments ...interface{}) string {
	operation := fmt.Sprintf(format, arguments...)

	localIdentifier := fragment.identifierTracker.NewIdentifier()

	instruction := &Instruction{Result: localIdentifier, Operation: operation}

	fragment.Lines = append(fragment.Lines, instruction)
	fragment.result = localIdentifier

	return localIdentifier
}

// AddLabel adds one single label to the fragment
func (fragment *Fragment) AddLabel(label string) {
	fragment.Lines = append(fragment.Lines, &Label{label})
}

func (fragment *Fragment) GetResult() string {
	if fragment.ResultType == VOID {
		panic("fragment generates void, shouldn't produce any local variable")
	}
	if fragment.result == "" {
		panic("last instruction does not produce any result")
	}
	return fragment.result
}

func (fragment *Fragment) String() string {
	var result string

	for _, line := range fragment.Lines {
		result += line.Line()
		result += "\n"
	}

	return result
}
