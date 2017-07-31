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

// Next returns the next valid token, or ILLEGAL with none-nil error
func (scanner *Scanner) Next() (*token.Token, error) {
	scanner.skipWhitespaces()

	if !scanner.file.IsEOF() {
		ch := scanner.file.NextChar()

		switch {
		case isDigit(ch):
			scanner.cur = string(ch)
			break
		default:
		}
	}

	return token.EOFToken(), nil
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
