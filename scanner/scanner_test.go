package scanner

import (
	"fmt"
	"testing"

	"github.com/carlcui/expressive/input"
	"github.com/carlcui/expressive/token"
)

func TestScanOperator(t *testing.T) {
	for operatorRaw, operatorType := range token.GetOperatorsMapping() {
		testScanningOneToken(operatorRaw, operatorRaw, operatorType, t)
	}
}

func TestScanKeyword(t *testing.T) {
	for keywordRaw, keywordType := range token.GetKeywordsMapping() {
		testScanningOneToken(keywordRaw, keywordRaw, keywordType, t)
	}
}

func TestScanInteger(t *testing.T) {
	actuals := []string{
		"123",
		"456",
		"0",
		"-5",
		"-0",
	}

	for _, actual := range actuals {
		testScanningOneToken(actual, actual, token.INT_LITERAL, t)
	}
}

func TestScanFloat(t *testing.T) {
	actuals := []string{
		"123.123",
		"123.",
		"0.3",
		"000.2",
		"0123.5",
		"-5.5",
		"-0.123",
	}

	for _, actual := range actuals {
		testScanningOneToken(actual, actual, token.FLOAT_LITERAL, t)
	}
}

func TestScanStringLiteralSuccess(t *testing.T) {
	actuals := []string{
		"\"abc\"",
		"\"\"",
		"\"  \"",
		"\"\\\"\"",
		"\"\\'\"",
		"\"\\n\"",
		"\"\\t\"",
		"\"\\0\"",
		"\"\\\\\"",
	}

	for _, actual := range actuals {
		testScanningOneToken(actual, actual, token.STRING_LITERAL, t)
	}
}

func TestScanCharacterLiteralSuccess(t *testing.T) {
	actuals := []string{
		"'a'",
		"'1'",
		"'🙃'",
		"'你'",
		"'\\0'",
		"'\\''",
		"'\\\"'",
		"'\\t'",
		"'\\n'",
		"'\\\\'",
	}

	for _, actual := range actuals {
		testScanningOneToken(actual, actual, token.CHAR_LITERAL, t)
	}
}

func TestScanSingleComment(t *testing.T) {
	stringInputs := []string{
		"// abc",
		"   // deo //",
		"// ddd*/",
		"// feeffe \n",
		"/* abc */",
		"/* // */",
		"/* *** */",
		"/* /n /n // */",
		"/* /n */ /n",
	}

	expected := []string{
		"// abc",
		"// deo //",
		"// ddd*/",
		"// feeffe ",
		"/* abc */",
		"/* // */",
		"/* *** */",
		"/* /n /n // */",
		"/* /n */",
	}

	for i, input := range stringInputs {
		testScanningOneToken(input, expected[i], token.COMMENT, t)
	}
}

func testScanningOneToken(stringInput string, expected string, tokenType token.Type, t *testing.T) {
	var input input.StringInput

	input.Init(stringInput)

	var scanner ExpressiveScanner
	scanner.Init(&input)

	tok := scanner.Next()

	expectedToken := token.Token{TokenType: tokenType, Raw: expected}

	compareTokens(*tok, expectedToken, t)
}

func testScanningTokens(stringInput string, expectedTokens []token.Token, t *testing.T) {
	var input input.StringInput

	input.Init(stringInput)

	var scanner ExpressiveScanner
	scanner.Init(&input)

	for _, expected := range expectedTokens {
		tok := scanner.Next()

		compareTokens(*tok, expected, t)
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
