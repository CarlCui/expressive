package token

import "strconv"
import "strings"

// Type is the type of a token
type Type int

const (
	// ILLEGAL is a type of token that is not specified in language spec (should not exist)
	ILLEGAL Type = iota
	EOF
	COMMENT

	// data types
	LITERAL
	IDENTIFIER
	INT
	FLOAT
	CHAR
	BOOLEAN
	STRING

	operatorStart
	// operators
	ADD
	SUB
	MUL
	DIV

	LAND
	LOR
	LNOT

	ASSIGN
	EQUAL

	SEMI // SEMI: semi-colon (;)
	COLON
	operatorEnd

	keywordStart
	// keywords
	LET
	CONST

	PRINT
	keywordEnd
)

// A Token represents a mainingful word in a program.
type Token struct {
	TokenType Type
	Raw       string
}

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	LITERAL:    "LITERAL",
	IDENTIFIER: "IDENTIFIER",
	INT:        "INT",
	FLOAT:      "FLOAT",
	CHAR:       "CHAR",
	BOOLEAN:    "BOOLEAN",
	STRING:     "STRING",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	LAND: "&&", // logic and
	LOR:  "||", // logic or
	LNOT: "!",  // logic not

	ASSIGN: "=",
	EQUAL:  "==",

	SEMI:  ";",
	COLON: ":",

	LET:   "let",
	CONST: "const",

	PRINT: "print",
}

var keywords map[string]Type
var operators map[string]Type

func init() {
	// init keywords mapping
	keywords = make(map[string]Type)
	for i := keywordStart + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = i
	}

	operators = make(map[string]Type)
	for i := operatorStart + 1; i < operatorEnd; i++ {
		operators[tokens[i]] = i
	}

}

func (tok *Token) String() string {
	if tok.TokenType >= 0 && int(tok.TokenType) < len(tokens) {
		return tokens[tok.TokenType] + ": " + tok.Raw
	}

	return "token(" + strconv.Itoa(int(tok.TokenType)) + ")"
}

// IllegalToken is a factory for generating a default illegal token
func IllegalToken(raw string) *Token {
	return &Token{TokenType: ILLEGAL, Raw: raw}
}

// EOFToken is a factory for generating a default EOF token
func EOFToken() *Token {
	return &Token{TokenType: EOF, Raw: ""}
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
