package codegen

import (
	"strconv"
	"strings"
)

// Labeller manages creating new labels that are never used before
// state is managed on a module level
type Labeller struct {
	index int
}

// Reset the label index to 0
func (labeller *Labeller) Reset() {
	labeller.index = 0
}

// NewSet increments the label index, making sure the label is not used before
func (labeller *Labeller) NewSet(tags ...string) string {
	labeller.index++

	return labeller.Label(tags...)
}

// Label creates a new label under the current set
func (labeller *Labeller) Label(tags ...string) string {
	return strings.Join(tags, ".") + strconv.Itoa(labeller.index)
}
