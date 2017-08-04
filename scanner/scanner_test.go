package scanner

import (
	"fmt"
	"testing"

	"github.com/carlcui/expressive/file"
	"github.com/carlcui/expressive/token"
)

func TestScanTokens(t *testing.T) {
	testFileDir := "."
	testFileName := "tokens.txt"

	var file file.File
	file.Init(testFileDir, testFileName)

	var scanner Scanner
	scanner.Init(&file)

	var expected = []token.Token{
		token.Token{TokenType: token.INT, Raw: "123"},
		token.Token{TokenType: token.FLOAT, Raw: "12.1"},
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
		tok, err := scanner.Next()

		if err != nil {
			t.Error(err)
		}

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
