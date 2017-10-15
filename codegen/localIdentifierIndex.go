package codegen

import (
	"strconv"
)

// LocalIdentifierTracker manages creating new local variables in current function scope
// State is managed on a function level
type LocalIdentifierTracker struct {
	index int
}

func (localVariableIndex *LocalIdentifierTracker) Reset() {
	localVariableIndex.index = 0
}

func (localVariableIndex *LocalIdentifierTracker) LocalVariable() string {
	return AsLocalVariable("l" + strconv.Itoa(localVariableIndex.index))
}

func (localVariableIndex *LocalIdentifierTracker) NewLocalVariable() string {
	localVariableIndex.index++

	return localVariableIndex.LocalVariable()
}

func AsLocalVariable(identifier string) string {
	return "%" + identifier
}
