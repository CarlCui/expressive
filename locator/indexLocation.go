package locator

import (
	"strconv"
)

type IndexLocation struct {
	Index int
}

func (loc *IndexLocation) Locate() string {
	return "at " + strconv.Itoa(loc.Index)
}

func CreateIndexLocation(index int) *IndexLocation {
	var loc IndexLocation
	loc.Index = index

	return &loc
}
