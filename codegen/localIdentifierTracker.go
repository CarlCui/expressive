package codegen

import (
	"strconv"
)

// LocalIdentifierTracker manages creating new local variables in current function scope
// State is managed on a function level
type LocalIdentifierTracker struct {
	index int
}

func (localIdentifierTracker *LocalIdentifierTracker) Reset() {
	localIdentifierTracker.index = 0
}

func (localIdentifierTracker *LocalIdentifierTracker) CurrentIdentifier() string {
	return AsLocalVariable("l" + strconv.Itoa(localIdentifierTracker.index))
}

func (localIdentifierTracker *LocalIdentifierTracker) NewIdentifier() string {
	localIdentifierTracker.index++

	return localIdentifierTracker.CurrentIdentifier()
}

func AsLocalVariable(identifier string) string {
	return "%" + identifier
}
