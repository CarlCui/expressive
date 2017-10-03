package token

import "strings"

import "github.com/carlcui/expressive/locator"

// A Token represents a mainingful word in a program.
type Token struct {
	TokenType Type
	Raw       string
	Locator   locator.Locator
}

func (tok *Token) String() string {
	return tok.TokenType.String() + ": " + tok.Raw
}

func (tok *Token) GetLocation() string {
	if tok.Locator == nil {
		return "unknown location"
	}

	return tok.Locator.Locate()
}

// IllegalToken is a factory for generating a default illegal token
func IllegalToken(raw string, locator locator.Locator) *Token {
	return &Token{TokenType: ILLEGAL, Raw: raw, Locator: locator}
}

// EOFToken is a factory for generating a default EOF token
func EOFToken(locator locator.Locator) *Token {
	return &Token{TokenType: EOF, Raw: "", Locator: locator}
}

// MatchKeyword returns a token containing info about that keyword, or nil if
// input is not a keyword
func MatchKeyword(reading string) *Token {
	if tokenType, isKeyword := keywords[reading]; isKeyword {
		return &Token{TokenType: tokenType, Raw: reading}
	}

	return nil
}

func MatchOperator(reading string) *Token {
	if tokenType, isOperator := operators[reading]; isOperator {
		return &Token{TokenType: tokenType, Raw: reading}
	}

	return nil
}

func HasOperatorPrefix(reading string) bool {
	for i := operatorStart + 1; i < operatorEnd; i++ {
		if strings.HasPrefix(tokens[i], reading) {
			return true
		}
	}
	return false
}
