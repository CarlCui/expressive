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
	VOID_OPERATOR
	ERROR_OPERATOR
)

var operatorStrings = [...]string{
	ADD:               "+(Addition)",
	SUBTRACT:          "-(Subtraction)",
	MULTIPLY:          "*(Multiplication)",
	DIVIDE:            "/(Division)",
	MODULO:            "%(Modulus)",
	EXPONENTIATE:      "^^(Exponentiation)",
	LOGIC_AND:         "&&(Logical and)",
	LOGIC_NOT:         "!(Logical not)",
	LOGIC_OR:          "||(Logical or)",
	IF_ELSE:           "? :(Ternary if else)",
	GREATER:           ">(Greater)",
	GREATER_OR_EQUAL:  ">=(Greater or equal)",
	LESS:              "<(Less)",
	LESS_OR_EQUAL:     "<=(Less or equal)",
	SHALLOW_EQUAL:     "==(Shallow equal)",
	SHALLOW_NOT_EQUAL: "!=(Shallow not euqal)",
	DEEP_EQUAL:        "===(Deep equal)",
	DEEP_NOT_EQUAL:    "!== (Deep not equal)",
	VOID_OPERATOR:     "void",
	ERROR_OPERATOR:    "??? (Error operator)",
}

func (operator Operator) String() string {
	if operator >= ADD && operator <= ERROR_OPERATOR {
		return operatorStrings[operator]
	}

	return operatorStrings[ERROR_OPERATOR]
}

func GetOperator(tok *token.Token) Operator {
	switch tok.TokenType {
	case token.ADD, token.ASSIGN_ADD:
		return ADD
	case token.SUB, token.ASSIGN_SUB:
		return SUBTRACT
	case token.MUL, token.ASSIGN_MUL:
		return MULTIPLY
	case token.DIV, token.ASSIGN_DIV:
		return DIVIDE
	case token.MOD, token.ASSIGN_MOD:
		return MODULO
	case token.POW, token.ASSIGN_POW:
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
