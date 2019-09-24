package token

import "strconv"

// Type is the type of a token
type Type int

const (
	// ILLEGAL is a type of token that is not specified in language spec (should not exist)
	ILLEGAL Type = iota
	EOF
	COMMENT

	LITERAL
	IDENTIFIER

	// consts
	INT_LITERAL
	FLOAT_LITERAL
	CHAR_LITERAL
	BOOLEAN_LITERAL
	STRING_LITERAL

	operatorStart
	// operators
	ADD
	SUB
	MUL
	DIV
	MOD
	POW

	LAND
	LOR
	LNOT

	ASSIGN
	ASSIGN_ADD
	ASSIGN_SUB
	ASSIGN_MUL
	ASSIGN_DIV
	ASSIGN_MOD
	ASSIGN_POW

	EQUAL
	NOT_EQUAL // not equal
	LESS
	LEQ // less or equal
	GREATER
	GEQ              // greater or equal
	TRIPLE_EQUAL     // deep equal
	TRIPLE_NOT_EQUAL // deep not equal

	LEFT_PAREN
	RIGHT_PAREN

	LEFT_CURLY_BRACE
	RIGHT_CURLY_BRACE

	QUESTION_MARK

	SEMI // SEMI: semi-colon (;)
	COLON
	COMMA
	operatorEnd

	keywordStart
	// keywords
	LET
	CONST

	IF
	ELSE

	WHILE
	FOR
	IN
	BREAK

	SWITCH
	CASE
	DEFAULT

	INT_KEYWORD
	FLOAT_KEYWORD
	CHAR_KEYWORD
	BOOL_KEYWORD
	STRING_KEYWORD

	TRUE
	FALSE

	PRINT
	keywordEnd
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	LITERAL:         "LITERAL",
	IDENTIFIER:      "IDENTIFIER",
	INT_LITERAL:     "INTLITERAL",
	FLOAT_LITERAL:   "FLOATLITERAL",
	CHAR_LITERAL:    "CHARLITERAL",
	BOOLEAN_LITERAL: "BOOLEANLITERAL",
	STRING_LITERAL:  "STRINGLITERAL",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",
	MOD: "%",
	POW: "^^",

	LAND: "&&", // logic and
	LOR:  "||", // logic or
	LNOT: "!",  // logic not

	ASSIGN:     "=",
	ASSIGN_ADD: "+=",
	ASSIGN_SUB: "-=",
	ASSIGN_MUL: "*=",
	ASSIGN_DIV: "/=",
	ASSIGN_MOD: "%=",
	ASSIGN_POW: "^^=",

	EQUAL:            "==",
	NOT_EQUAL:        "!=",
	LESS:             "<",
	LEQ:              "<=",
	GREATER:          ">",
	GEQ:              ">=",
	TRIPLE_EQUAL:     "===",
	TRIPLE_NOT_EQUAL: "!==",

	LEFT_PAREN:  "(",
	RIGHT_PAREN: ")",

	LEFT_CURLY_BRACE:  "{",
	RIGHT_CURLY_BRACE: "}",

	QUESTION_MARK: "?",

	SEMI:  ";",
	COLON: ":",
	COMMA: ",",

	LET:   "let",
	CONST: "const",

	IF:   "if",
	ELSE: "else",

	WHILE: "while",
	FOR:   "for",
	IN:    "in",
	BREAK: "break",

	SWITCH:  "switch",
	CASE:    "case",
	DEFAULT: "default",

	INT_KEYWORD:    "int",
	FLOAT_KEYWORD:  "float",
	CHAR_KEYWORD:   "char",
	BOOL_KEYWORD:   "bool",
	STRING_KEYWORD: "string",

	TRUE:  "true",
	FALSE: "false",

	PRINT: "print",
}

func (tokenType Type) String() string {
	if tokenType >= 0 && int(tokenType) < len(tokens) {
		return tokens[tokenType]
	}

	return "token(" + strconv.Itoa(int(tokenType)) + ")"
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

func GetKeywordsMapping() map[string]Type {
	return keywords
}

func GetOperatorsMapping() map[string]Type {
	return operators
}
