package signature

import "github.com/carlcui/expressive/token"

type Operator int

const (
	ADD Operator = iota
	SUBTRACT
	MULTIPLY
	DIVIDE
	MODULO
	EXPONENTIATE
	LOGIC_AND
	LOGIC_NOT
	LOGIC_OR
	IF_ELSE // ? :
	GREATER
	GREATER_OR_EQUAL
	LESS
	LESS_OR_EQUAL
	SHALLOW_EQUAL
	DEEP_EQUAL
	SHALLOW_NOT_EQUAL
	DEEP_NOT_EQUAL
	ERROR_OPERATOR
)

func GetOperator(tok *token.Token) Operator {
	switch tok.TokenType {
	case token.ADD:
		return ADD
	case token.SUB:
		return SUBTRACT
	case token.MUL:
		return MULTIPLY
	case token.DIV:
		return DIVIDE
	case token.MOD:
		return MODULO
	case token.POW:
		return EXPONENTIATE
	case token.LAND:
		return LOGIC_AND
	case token.LNOT:
		return LOGIC_NOT
	case token.LOR:
		return LOGIC_OR
	case token.GREATER:
		return GREATER
	case token.GEQ:
		return GREATER_OR_EQUAL
	case token.LESS:
		return LESS
	case token.LEQ:
		return LESS_OR_EQUAL
	case token.EQUAL:
		return SHALLOW_EQUAL
	case token.NOT_EQUAL:
		return SHALLOW_NOT_EQUAL
	case token.TRIPLE_EQUAL:
		return DEEP_EQUAL
	case token.TRIPLE_NOT_EQUAL:
		return DEEP_NOT_EQUAL
	default:
		return ERROR_OPERATOR
	}
}
