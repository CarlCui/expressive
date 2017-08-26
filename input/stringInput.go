package input

import (
	"strconv"
	"unicode/utf8"

	"github.com/carlcui/expressive/locator"
)

type StringInput struct {
	curPos int

	src []byte
}

func (input *StringInput) Init(src string) {
	input.curPos = 0
	input.src = []byte(src)
}

func (input *StringInput) NextChar() rune {
	if input.IsEOF() {
		panic("eof at " + strconv.Itoa(input.curPos))
	}

	r, size := utf8.DecodeRune(input.src[input.curPos:])

	if r == utf8.RuneError {
		panic("error decoding rune at " + strconv.Itoa(input.curPos))
	}

	input.curPos += size

	return r
}

func (input *StringInput) Peek() rune {
	if input.IsEOF() {
		panic("eof at " + strconv.Itoa(input.curPos))
	}

	r, _ := utf8.DecodeRune(input.src[input.curPos:])

	if r == utf8.RuneError {
		panic("error decoding rune at " + strconv.Itoa(input.curPos))
	}

	return r
}

func (input *StringInput) IsEOF() bool {
	return input.curPos >= len(input.src)
}

func (input *StringInput) CurLoc() locator.Locator {
	var loc locator.IndexLocation
	loc.Index = input.curPos

	return &loc
}
