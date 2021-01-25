package scanner

import (
	"github.com/carlcui/expressive/locator"
	"github.com/carlcui/expressive/token"
)

type MockScanner struct {
	toks []*token.Token
	pos  int
}

func (scanner *MockScanner) Init(toks []*token.Token) {
	scanner.toks = toks
	scanner.pos = 0
}

func (scanner *MockScanner) Next() *token.Token {
	if scanner.pos >= len(scanner.toks) {
		return token.IllegalToken("", locator.CreateIndexLocation(scanner.pos))
	}

	cur := scanner.toks[scanner.pos]
	cur.Locator = &locator.IndexLocation{Index: scanner.pos}

	scanner.pos++

	return cur
}
