package token

import "strconv"

// Type is the type of a token
type Type int

const (
	// ILLEGAL is a type of token that is not specified in language spec (should not exist)
	ILLEGAL Type = iota
	EOF
	COMMENT

	// data types
	IDENTIFIER
	INT
	FLOAT
	CHAR
	BOOLEAN
	STRING

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

	// keywords
	LET
	CONST
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

	SEMI: ";",

	LET:   "let",
	CONST: "const",
}

var keywords map[string]Type
var operators map[string]Type

func init() {
	// init keywords mapping
	keywords = make(map[string]Type)
	for i := LET; i < CONST; i++ {
		keywords[tokens[i]] = i
	}

	operators = make(map[string]Type)
	for i := ADD; i < SEMI; i++ {
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

// IsKeyword returns a token containing info about that keyword, or nil if
// input is not a keyword
func IsKeyword(reading string) *Token {
	if tokenType, isKeyword := keywords[reading]; isKeyword {
		return &Token{TokenType: tokenType, Raw: reading}
	}

	return nil
}
