package scanner

import (
	"fmt"
	"testing"

	"github.com/carlcui/expressive/input"
	"github.com/carlcui/expressive/token"
)

func TestScanTokens(t *testing.T) {
	testFileDir := "."
	testFileName := "tokens.txt"

	var file input.File
	file.Init(testFileDir, testFileName)

	var scanner ExpressiveScanner
	scanner.Init(&file)

	var expected = []token.Token{
		token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		token.Token{TokenType: token.FLOAT_LITERAL, Raw: "12.1"},
		token.Token{TokenType: token.IDENTIFIER, Raw: "abc"},
		token.Token{TokenType: token.LET, Raw: "let"},
		token.Token{TokenType: token.ADD, Raw: "+"},
		token.Token{TokenType: token.SUB, Raw: "-"},
		token.Token{TokenType: token.MUL, Raw: "*"},
		token.Token{TokenType: token.DIV, Raw: "/"},
		token.Token{TokenType: token.LAND, Raw: "&&"},
		token.Token{TokenType: token.LOR, Raw: "||"},
		token.Token{TokenType: token.LNOT, Raw: "!"},
	}

	for _, expectedToken := range expected {
		tok := scanner.Next()

		compareTokens(*tok, expectedToken, t)
	}

}

func TestScanInteger(t *testing.T) {
	var input input.StringInput

	input.Init("123 456 0")

	var scanner ExpressiveScanner
	scanner.Init(&input)

	var expected = []token.Token{
		token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		token.Token{TokenType: token.INT_LITERAL, Raw: "456"},
		token.Token{TokenType: token.INT_LITERAL, Raw: "0"},
		token.Token{TokenType: token.EOF, Raw: ""},
	}

	for _, expectedToken := range expected {
		tok := scanner.Next()

		compareTokens(*tok, expectedToken, t)
	}
}

func TestScanFloat(t *testing.T) {
	var input input.StringInput

	input.Init("123.123 123. 0.3 000.2 0123.5")

	var scanner ExpressiveScanner
	scanner.Init(&input)

	var expected = []token.Token{
		token.Token{TokenType: token.FLOAT_LITERAL, Raw: "123.123"},
		token.Token{TokenType: token.FLOAT_LITERAL, Raw: "123."},
		token.Token{TokenType: token.FLOAT_LITERAL, Raw: "0.3"},
		token.Token{TokenType: token.FLOAT_LITERAL, Raw: "000.2"},
		token.Token{TokenType: token.FLOAT_LITERAL, Raw: "0123.5"},
		token.Token{TokenType: token.EOF, Raw: ""},
	}

	for _, expectedToken := range expected {
		tok := scanner.Next()

		compareTokens(*tok, expectedToken, t)
	}
}

func TestScanStringLiteralSuccess(t *testing.T) {
	var input input.StringInput

	input.Init("\"abc\" \"\" \"  \" \"\\\"\" \"\\'\" \"\\n\" \"\\t\" \"\\0\" \"\\\\\"  ")

	var scanner ExpressiveScanner
	scanner.Init(&input)

	var expected = []token.Token{
		token.Token{TokenType: token.STRING_LITERAL, Raw: "\"abc\""},
		token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\""},
		token.Token{TokenType: token.STRING_LITERAL, Raw: "\"  \""},
		token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\\"\""},
		token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\'\""},
		token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\n\""},
		token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\t\""},
		token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\0\""},
		token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\\\\""},
	}

	for _, expectedToken := range expected {
		tok := scanner.Next()

		compareTokens(*tok, expectedToken, t)
	}
}

func compareTokens(actual token.Token, expected token.Token, t *testing.T) {
	equal := actual.TokenType == expected.TokenType && actual.Raw == expected.Raw

	if !equal {
		reportTokenParsingError(actual, expected, t)
	} else {
		fmt.Println(actual)
	}

}

func reportTokenParsingError(actual token.Token, expected token.Token, t *testing.T) {
	t.Errorf("Actual: %s, expected: %s \n", actual.String(), expected.String())
}
