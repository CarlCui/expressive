package token

import "strconv"

type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF
	COMMENT

	IDENTIFIER
	INT
	FLOAT
	CHAR
	BOOLEAN
	STRING

	ADD
	SUB
	MUL
	DIV

	AND
	OR
	NOT
)

type Token struct {
	TokenType TokenType
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

	ADD: "ADD",
	SUB: "SUB",
	MUL: "MUL",
	DIV: "DIV",

	AND: "AND",
	OR:  "OR",
	NOT: "NOT",
}

func (tok *Token) String() string {
	if tok.TokenType >= 0 && int(tok.TokenType) < len(tokens) {
		return tokens[tok.TokenType] + ": " + tok.Raw
	}

	return "token(" + strconv.Itoa(int(tok.TokenType)) + ")"
}

// IllegalToken is a factory for generating a default illegal token
func IllegalToken() *Token {
	return &Token{TokenType: ILLEGAL, Raw: ""}
}
