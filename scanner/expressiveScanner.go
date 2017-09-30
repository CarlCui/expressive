package scanner

import "github.com/carlcui/expressive/token"
import "github.com/carlcui/expressive/input"
import "github.com/carlcui/expressive/locator"
import "unicode"

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
		ch := scanner.input.NextChar()

		scanner.cur = string(ch)

		switch {
		case isDigit(ch):
			tok = scanner.parseNumber()
			break
		case isIdentifierStart(ch):
			tok = scanner.parseIdentifier()
			break
		case isDoubleQuote(ch):
			tok = scanner.parseStringLiteral()
			break
		case token.HasOperatorPrefix(string(ch)):
			tok = scanner.parseOperator()
			break
		default:
			tok = token.IllegalToken(scanner.cur, scanner.curLoc)
		}
	}

	return tok
}

/*
	a number is either:

	1. an integer: [0-9]+
	2. a float: [0-9]+.[0-9]+

*/
func (scanner *ExpressiveScanner) parseNumber() *token.Token {
	// parse int
	scanner.appendConsequentDigits()

	// parse float
	if !scanner.input.IsEOF() && isDot(scanner.input.Peek()) {
		scanner.cur += string(scanner.input.NextChar())

		scanner.appendConsequentDigits()

		return &token.Token{TokenType: token.FLOAT_LITERAL, Raw: scanner.cur, Locator: scanner.curLoc}
	}

	return &token.Token{TokenType: token.INT_LITERAL, Raw: scanner.cur, Locator: scanner.curLoc}
}

func (scanner *ExpressiveScanner) appendConsequentDigits() {
	for !scanner.input.IsEOF() && isDigit(scanner.input.Peek()) {
		ch := scanner.input.NextChar()

		scanner.cur += string(ch)
	}
}

/*
	identifier := [_a-zA-Z][_0-9a-zA-Z]*
*/
func (scanner *ExpressiveScanner) parseIdentifier() *token.Token {
	for !scanner.input.IsEOF() && !isWhitespace(scanner.input.Peek()) && (isIdentifierStart(scanner.input.Peek()) || isDigit(scanner.input.Peek())) {
		scanner.cur += string(scanner.input.NextChar())
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

		cur := scanner.input.NextChar()
		scanner.cur += string(cur)

		if isBackSlash(cur) {
			// parse escape sequence

			if scanner.input.IsEOF() {
				return token.IllegalToken(scanner.cur, loc)
			}

			if !isControlSequenceCharacter(scanner.input.Peek()) {
				return token.IllegalToken(scanner.cur, loc)
			}

			scanner.cur += string(scanner.input.NextChar())

		}
	}

	// expecting back qoute
	if scanner.input.IsEOF() {
		return token.IllegalToken(scanner.cur, loc)
	}
	if isReturn(scanner.input.Peek()) {
		return token.IllegalToken(scanner.cur, loc)
	}

	if !isDoubleQuote(scanner.input.Peek()) {
		return token.IllegalToken(scanner.cur, loc)
	}

	scanner.cur += string(scanner.input.NextChar()) // "

	return &token.Token{TokenType: token.STRING_LITERAL, Raw: scanner.cur, Locator: loc}
}

/*
	charLiteral := '[^'^\n]|(\asciiEscapeControl)'
*/

/*
	operators
*/
func (scanner *ExpressiveScanner) parseOperator() *token.Token {
	for !scanner.input.IsEOF() && token.HasOperatorPrefix(scanner.cur+string(scanner.input.Peek())) {
		scanner.cur += string(scanner.input.NextChar())
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
		scanner.input.NextChar()
	}
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

func isDoubleQuote(ch rune) bool {
	return ch == '"'
}

func isSingleQuote(ch rune) bool {
	return ch == '\''
}

func isBackSlash(ch rune) bool {
	return ch == '\\'
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
