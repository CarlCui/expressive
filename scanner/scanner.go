package scanner

import "github.com/carlcui/expressive/token"
import "github.com/carlcui/expressive/file"
import "unicode"

// Scanner is a lexical analyzer for expressive
type Scanner struct {
	file *file.File

	cur string // current string buffer
}

// Init initializes scanner, setting current string buffer to empty string
func (scanner *Scanner) Init(file *file.File) {
	scanner.file = file
	scanner.cur = ""
}

// Next returns the next valid token, or ILLEGAL if parsing failed
func (scanner *Scanner) Next() *token.Token {
	scanner.skipWhitespaces()

	tok := token.EOFToken()

	if !scanner.file.IsEOF() {
		ch := scanner.file.NextChar()

		scanner.cur = string(ch)

		switch {
		case isDigit(ch):
			tok = scanner.parseNumber()
			break
		case isIdentifierStart(ch):
			tok = scanner.parseIdentifier()
			break
		case token.HasOperatorPrefix(string(ch)):
			tok = scanner.parseOperator()
			break
		default:
			tok = token.IllegalToken(scanner.cur)
		}
	}

	return tok
}

/*
	a number is either:

	1. an integer: [0-9]+
	2. a float: [0-9]+.[0-9]+

*/
func (scanner *Scanner) parseNumber() *token.Token {
	// parse int
	scanner.appendConsequentDigits()

	// parse float
	if !scanner.file.IsEOF() && isDot(scanner.file.Peek()) {
		scanner.cur += string(scanner.file.NextChar())

		scanner.appendConsequentDigits()

		return &token.Token{TokenType: token.FLOAT, Raw: scanner.cur}
	}

	return &token.Token{TokenType: token.INT, Raw: scanner.cur}
}

func (scanner *Scanner) appendConsequentDigits() {
	for !scanner.file.IsEOF() && isDigit(scanner.file.Peek()) {
		ch := scanner.file.NextChar()

		scanner.cur += string(ch)
	}
}

/*
	identifier := [_a-zA-Z][_0-9a-zA-Z]*

*/
func (scanner *Scanner) parseIdentifier() *token.Token {
	for !scanner.file.IsEOF() && !isWhitespace(scanner.file.Peek()) && (isIdentifierStart(scanner.file.Peek()) || isDigit(scanner.file.Peek())) {
		scanner.cur += string(scanner.file.NextChar())
	}

	tok := token.MatchKeyword(scanner.cur)

	if tok == nil {
		tok = &token.Token{TokenType: token.IDENTIFIER, Raw: scanner.cur}
	}

	return tok
}

/*
	operators
*/
func (scanner *Scanner) parseOperator() *token.Token {
	for !scanner.file.IsEOF() && token.HasOperatorPrefix(scanner.cur+string(scanner.file.Peek())) {
		scanner.cur += string(scanner.file.NextChar())
	}

	tok := token.MatchOperator(scanner.cur)

	if tok == nil {
		tok = token.IllegalToken(scanner.cur)
	}

	return tok
}

func (scanner *Scanner) skipWhitespaces() {
	for !scanner.file.IsEOF() && isWhitespace(scanner.file.Peek()) {
		scanner.file.NextChar()
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

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

func isUnderscore(ch rune) bool {
	return ch == '_'
}

func isIdentifierStart(ch rune) bool {
	return isUnderscore(ch) || isLetter(ch)
}
