package input

import "github.com/carlcui/expressive/locator"

type Input interface {
	NextChar() rune
	Peek() rune
	IsEOF() bool
	CurLoc() locator.Locator
}
