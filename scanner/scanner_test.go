package scanner

import (
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
	}
}

func reportTokenParsingError(actual token.Token, expected token.Token, t *testing.T) {
	t.Errorf("Actual: %s, expected: %s \n", actual.String(), expected.String())
}
