package input

import "github.com/carlcui/expressive/locator"

type Input interface {
	NextRune() rune
	Peek() rune
	IsEOF() bool
	CurLoc() locator.Locator
}
