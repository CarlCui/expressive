package token

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

	QUESTION_MARK

	SEMI // SEMI: semi-colon (;)
	COLON
	operatorEnd

	keywordStart
	// keywords
	LET
	CONST

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

	ASSIGN: "=",

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

	QUESTION_MARK: "?",

	SEMI:  ";",
	COLON: ":",

	LET:   "let",
	CONST: "const",

	INT_KEYWORD:    "int",
	FLOAT_KEYWORD:  "float",
	CHAR_KEYWORD:   "char",
	BOOL_KEYWORD:   "bool",
	STRING_KEYWORD: "string",

	TRUE:  "true",
	FALSE: "false",

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

func GetKeywordsMapping() map[string]Type {
	return keywords
}

func GetOperatorsMapping() map[string]Type {
	return operators
}
