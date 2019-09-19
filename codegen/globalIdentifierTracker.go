package codegen

import (
	"strconv"
)

// GlobalIdentifierTracker manages creating new global identifiers
// State is managed on a function level
type GlobalIdentifierTracker struct {
	index int
}

func (globalIdentifierTracker *GlobalIdentifierTracker) Reset() {
	globalIdentifierTracker.index = 0
}

func (globalIdentifierTracker *GlobalIdentifierTracker) CurrentIdentifier() string {
	return "g" + strconv.Itoa(globalIdentifierTracker.index)
}

func (globalIdentifierTracker *GlobalIdentifierTracker) NewIdentifier() string {
	globalIdentifierTracker.index++

	return globalIdentifierTracker.CurrentIdentifier()
}
