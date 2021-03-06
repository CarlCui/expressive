package scanner

import (
	"unicode"

	"github.com/carlcui/expressive/input"
	"github.com/carlcui/expressive/locator"
	"github.com/carlcui/expressive/token"
)

// ExpressiveScanner is a lexical analyzer for expressive
type ExpressiveScanner struct {
	input input.Input

	cur    string          // current string buffer
	curLoc locator.Locator // current location
}

// Init initializes scanner, setting current string buffer to empty string
func (scanner *ExpressiveScanner) Init(input input.Input) {
	scanner.input = input
	scanner.cur = ""
}

// Next returns the next valid token, or ILLEGAL if parsing failed
func (scanner *ExpressiveScanner) Next() *token.Token {
	scanner.skipWhitespaces()

	scanner.curLoc = scanner.input.CurLoc()

	tok := token.EOFToken(scanner.curLoc)

	if !scanner.input.IsEOF() {
		ch := scanner.input.NextRune()

		scanner.cur = string(ch)

		commentTok := scanner.tryParseComment(ch)

		if commentTok != nil {
			return commentTok
		}

		switch {
		case isDigit(ch) || isMinus(ch):
			tok = scanner.parseNumber(ch)
			break
		case isIdentifierStart(ch):
			tok = scanner.parseIdentifier()
			break
		case isDoubleQuote(ch):
			tok = scanner.parseStringLiteral()
			break
		case isSingleQuote(ch):
			tok = scanner.parseCharacterLiteral()
		case token.HasOperatorPrefix(string(ch)):
			tok = scanner.parseOperator()
			break
		default:
			tok = token.IllegalToken(scanner.cur, scanner.curLoc)
		}
	}

	return tok
}

func (scanner *ExpressiveScanner) tryParseComment(cur rune) *token.Token {
	if scanner.input.IsEOF() { // cannot be a comment
		return nil
	}

	if isSlash(cur) && isSlash(scanner.input.Peek()) { // `//`

		scanner.cur += string(scanner.input.NextRune())

		for !scanner.input.IsEOF() {
			if isReturn(scanner.input.Peek()) {
				return &token.Token{TokenType: token.COMMENT, Raw: scanner.cur, Locator: scanner.curLoc}
			}

			ch := scanner.input.NextRune()

			scanner.cur += string(ch)
		}

		return &token.Token{TokenType: token.COMMENT, Raw: scanner.cur, Locator: scanner.curLoc}

	} else if isSlash(cur) && isAsterisk(scanner.input.Peek()) { // `/*`

		for !scanner.input.IsEOF() {
			ch := scanner.input.NextRune()

			scanner.cur += string(ch)

			// */
			if isAsterisk(ch) && isSlash(scanner.input.Peek()) {
				scanner.cur += string(scanner.input.NextRune())
				return &token.Token{TokenType: token.COMMENT, Raw: scanner.cur, Locator: scanner.curLoc}
			}
		}
	}

	return nil
}

/*
	a number is either:

	1. an integer: -?[0-9]+
	2. a float: -?[0-9]+.[0-9]+

*/
func (scanner *ExpressiveScanner) parseNumber(first rune) *token.Token {

	if isMinus(first) {
		if scanner.input.IsEOF() || !isDigit(scanner.input.Peek()) { // cannot be a negative number
			return scanner.parseOperator()
		}
	}

	// parse int
	scanner.appendConsequentDigits()

	// parse float
	if !scanner.input.IsEOF() && isDot(scanner.input.Peek()) {
		scanner.cur += string(scanner.input.NextRune())

		scanner.appendConsequentDigits()

		return &token.Token{TokenType: token.FLOAT_LITERAL, Raw: scanner.cur, Locator: scanner.curLoc}
	}

	return &token.Token{TokenType: token.INT_LITERAL, Raw: scanner.cur, Locator: scanner.curLoc}
}

func (scanner *ExpressiveScanner) appendConsequentDigits() {
	for !scanner.input.IsEOF() && isDigit(scanner.input.Peek()) {
		ch := scanner.input.NextRune()

		scanner.cur += string(ch)
	}
}

/*
	identifier := [_a-zA-Z][_0-9a-zA-Z]*
*/
func (scanner *ExpressiveScanner) parseIdentifier() *token.Token {
	for !scanner.input.IsEOF() && !isWhitespace(scanner.input.Peek()) && (isIdentifierStart(scanner.input.Peek()) || isDigit(scanner.input.Peek())) {
		scanner.cur += string(scanner.input.NextRune())
	}

	tok := token.MatchKeyword(scanner.cur)

	if tok == nil {
		tok = &token.Token{TokenType: token.IDENTIFIER, Raw: scanner.cur, Locator: scanner.curLoc}
	}

	tok.Locator = scanner.curLoc

	return tok
}

/*
	stringLiteral := "[^"^\n]*"
*/
func (scanner *ExpressiveScanner) parseStringLiteral() *token.Token {
	loc := scanner.curLoc

	for !scanner.input.IsEOF() && !isReturn(scanner.input.Peek()) && !isDoubleQuote(scanner.input.Peek()) {

		ok := scanner.tryParseEscapeControlSequence()

		if !ok {
			return token.IllegalToken(scanner.cur, loc)
		}
	}

	// expecting back qoute
	if scanner.input.IsEOF() || isReturn(scanner.input.Peek()) || !isDoubleQuote(scanner.input.Peek()) {
		return token.IllegalToken(scanner.cur, loc)
	}

	scanner.cur += string(scanner.input.NextRune()) // "

	return &token.Token{TokenType: token.STRING_LITERAL, Raw: scanner.cur, Locator: loc}
}

/*
	charLiteral := '[^'^\n]|(\asciiEscapeControl)'
*/
func (scanner *ExpressiveScanner) parseCharacterLiteral() *token.Token {
	loc := scanner.curLoc

	if scanner.input.IsEOF() || isReturn(scanner.input.Peek()) || isSingleQuote(scanner.input.Peek()) {
		return token.IllegalToken(scanner.cur, loc)
	}

	ok := scanner.tryParseEscapeControlSequence()

	if !ok {
		return token.IllegalToken(scanner.cur, loc)
	}

	if scanner.input.IsEOF() || !isSingleQuote(scanner.input.Peek()) {
		return token.IllegalToken(scanner.cur, loc)
	}

	scanner.cur += string(scanner.input.NextRune())

	return &token.Token{TokenType: token.CHAR_LITERAL, Raw: scanner.cur, Locator: loc}
}

/*
	operators
*/
func (scanner *ExpressiveScanner) parseOperator() *token.Token {
	for !scanner.input.IsEOF() && token.HasOperatorPrefix(scanner.cur+string(scanner.input.Peek())) {
		scanner.cur += string(scanner.input.NextRune())
	}

	tok := token.MatchOperator(scanner.cur)

	if tok == nil {
		tok = token.IllegalToken(scanner.cur, scanner.curLoc)
	}

	tok.Locator = scanner.curLoc

	return tok
}

func (scanner *ExpressiveScanner) skipWhitespaces() {
	for !scanner.input.IsEOF() && isWhitespace(scanner.input.Peek()) {
		scanner.input.NextRune()
	}
}

func (scanner *ExpressiveScanner) tryParseEscapeControlSequence() bool {
	cur := scanner.input.NextRune()
	scanner.cur += string(cur)

	if isBackSlash(cur) {
		if scanner.input.IsEOF() || !isControlSequenceCharacter(scanner.input.Peek()) {
			return false
		}

		scanner.cur += string(scanner.input.NextRune())
	}

	return true
}

// char helpers
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

func isWhitespace(ch rune) bool {
	return unicode.IsSpace(ch)
}

func isDot(ch rune) bool {
	return ch == '.'
}

func isMinus(ch rune) bool {
	return ch == '-'
}

func isDoubleQuote(ch rune) bool {
	return ch == '"'
}

func isSingleQuote(ch rune) bool {
	return ch == '\''
}

func isBackSlash(ch rune) bool {
	return ch == '\\'
}

func isSlash(ch rune) bool {
	return ch == '/'
}

func isAsterisk(ch rune) bool {
	return ch == '*'
}

// implemented partially
func isControlSequenceCharacter(ch rune) bool {
	return ch == 'n' || ch == '0' || ch == 't' || ch == '\'' || ch == '"' || ch == '\\'
}

func isReturn(ch rune) bool {
	return ch == '\n'
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isUnderscore(ch rune) bool {
	return ch == '_'
}

func isIdentifierStart(ch rune) bool {
	return isUnderscore(ch) || isLetter(ch)
}
