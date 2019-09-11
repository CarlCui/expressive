package parser

import (
	"testing"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/token"
)

func TestParseValidVariableDelarationStmt_LetStmt(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.LET, Raw: "let"},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "foo"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	root := parseWithMockTokens(toks, shouldHaveNoError(t))

	stmtNode := root.(*ast.ProgramNode).Chilren[0]

	if stmt, ok := stmtNode.(*ast.VariableDeclarationNode); ok {
		if stmt.Tok.TokenType != token.LET {
			t.Error()
		}

		if identifier, ok := stmt.Identifier.(*ast.IdentifierNode); ok {
			if identifier.Tok.Raw != "foo" {
				t.Error()
			}
		} else {
			t.Error()
		}

		if stmt.DeclaredType != nil || stmt.Expr != nil {
			t.Error()
		}
	} else {
		t.Error()
	}
}

func TestParseValidVariableDeclarationStmt_LetStmtWithDeclaredType(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "foo"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	root := parseWithMockTokens(toks, shouldHaveNoError(t))

	stmtNode := root.(*ast.ProgramNode).Chilren[0]

	if stmt, ok := stmtNode.(*ast.VariableDeclarationNode); ok {
		if stmt.Tok.TokenType != token.LET {
			t.Error()
		}

		if identifier, ok := stmt.Identifier.(*ast.IdentifierNode); ok {
			if identifier.Tok.Raw != "foo" {
				t.Error()
			}
		} else {
			t.Error()
		}

		if declaredType, ok := stmt.DeclaredType.(*ast.TypeLiteralNode); ok {
			if declaredType.Tok.TokenType != token.INT_KEYWORD {
				t.Error()
			}
		} else {
			t.Error()
		}

		if stmt.Expr != nil {
			t.Error()
		}
	} else {
		t.Error()
	}
}

func TestParseIntegerLiteral(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "0"},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "-1"},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "-0"},
	}

	expectedVals := []int{
		123,
		0,
		-1,
		0,
	}

	for index, tok := range toks {
		parser := initParserWithMockTokens([]*token.Token{tok})

		node := parser.parserInt()

		intNode, ok := node.(*ast.IntegerNode)

		if !ok || intNode.Val != expectedVals[index] {
			reportTestError("Error parsing integer literal node", node, t)
		}
	}
}

func TestParseFloatLiteral(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.FLOAT_LITERAL, Raw: "123.123"},
		&token.Token{TokenType: token.FLOAT_LITERAL, Raw: "123.0"},
		&token.Token{TokenType: token.FLOAT_LITERAL, Raw: "-1.0"},
		&token.Token{TokenType: token.FLOAT_LITERAL, Raw: "-123.5"},
		&token.Token{TokenType: token.FLOAT_LITERAL, Raw: "-0.6"},
		&token.Token{TokenType: token.FLOAT_LITERAL, Raw: "-0.0"},
	}

	expectedVals := []float32{
		123.123,
		123,
		-1,
		-123.5,
		-0.6,
		0,
	}

	for index, tok := range toks {
		parser := initParserWithMockTokens([]*token.Token{tok})

		node := parser.parseFloat()

		floatNode, ok := node.(*ast.FloatNode)

		if !ok || floatNode.Val != expectedVals[index] {
			reportTestError("Error parsing float literal node", node, t)
		}
	}
}

func TestParseBooleanLiteral(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.TRUE},
		&token.Token{TokenType: token.FALSE},
	}

	for _, tok := range toks {
		parser := initParserWithMockTokens([]*token.Token{tok})

		node := parser.parseBool()

		if _, ok := node.(*ast.BooleanNode); !ok {
			reportTestError("Error parsing boolean literal node", node, t)
		}
	}
}

func TestParseBooleanLiteralFailed(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.ADD},
		&token.Token{TokenType: token.IDENTIFIER},
	}

	for _, tok := range toks {
		parser := initParserWithMockTokens([]*token.Token{tok})

		node := parser.parseBool()

		if _, ok := node.(*ast.BooleanNode); ok {
			reportTestError("Parsing boolean literal should fail", node, t)
		}
	}
}

func TestParseCharacterLiteral(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.CHAR_LITERAL, Raw: "'a'"},
		&token.Token{TokenType: token.CHAR_LITERAL, Raw: "'1'"},
		&token.Token{TokenType: token.CHAR_LITERAL, Raw: "'\\0'"},
		&token.Token{TokenType: token.CHAR_LITERAL, Raw: "'\\''"},
		&token.Token{TokenType: token.CHAR_LITERAL, Raw: "'\\\"'"},
		&token.Token{TokenType: token.CHAR_LITERAL, Raw: "'\\t'"},
		&token.Token{TokenType: token.CHAR_LITERAL, Raw: "'\\n'"},
		&token.Token{TokenType: token.CHAR_LITERAL, Raw: "'\\\\'"},
	}

	expectedVals := []rune{
		'a',
		'1',
		'\x00',
		'\'',
		'"',
		'\t',
		'\n',
		'\\',
	}

	for index, tok := range toks {
		parser := initParserWithMockTokens([]*token.Token{tok})

		node := parser.parseCharacter()

		charNode, ok := node.(*ast.CharacterNode)

		if !ok || charNode.Val != expectedVals[index] {
			reportTestError("Error parsing character literal node", node, t)
		}
	}
}

func TestParseStringLiteral(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"abc\""},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\""},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"  \""},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\\"\""},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\'\""},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\n\""},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\t\""},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\0\""},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"\\\\\""},
	}

	expectedVals := []string{
		"abc",
		"",
		"  ",
		"\\\"",
		"\\'",
		"\\n",
		"\\t",
		"\\0",
		"\\\\",
	}

	for index, tok := range toks {
		parser := initParserWithMockTokens([]*token.Token{tok})

		node := parser.parseString()

		stringNode, ok := node.(*ast.StringNode)

		if !ok || stringNode.Val != expectedVals[index] {
			reportTestError("Error parsing string literal node", node, t)
		}
	}
}

func TestSkipCommentTokenWhenParsing_InFront(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.COMMENT},
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "foo"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestSkipCommentTokenWhenParsing_InBetween(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "foo"},
		&token.Token{TokenType: token.COMMENT},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.COMMENT},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestSkipCommentTokenWhenParsing_After(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "foo"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.COMMENT},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestSkipCommentTokenWhenParsing_AfterEOF(t *testing.T) {
	toks := []*token.Token{
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "foo"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
		&token.Token{TokenType: token.COMMENT},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestNoEOFShouldFail(t *testing.T) {
	// let bar: int; (no EOF)
	toks := []*token.Token{
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "bar"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.SEMI},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParsingPrintStmt(t *testing.T) {
	// print "123";
	toks := []*token.Token{
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingPrintStmtWithOneArg(t *testing.T) {
	// print "123", 123;
	toks := []*token.Token{
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingPrintStmtWithMultipleArgs(t *testing.T) {
	// print "123", 123, 123.5;
	toks := []*token.Token{
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.FLOAT_LITERAL, Raw: "123.5"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingPrintStmtWithMultipleArgs_Fail(t *testing.T) {
	// print "123", 123,;
	toks := []*token.Token{
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParsingIfStmtWithNoElse(t *testing.T) {
	/*
		if (true) {
			print "123", 123;
		}
	*/
	toks := []*token.Token{
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.TRUE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingIfStmtWithElse(t *testing.T) {
	/*
		if (true) {
			print "123", 123;
		} else {
			print "123", 123;
		}
	*/
	toks := []*token.Token{
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.TRUE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.ELSE},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingIfStmtWithIfElse(t *testing.T) {
	/*
		if (true) {
			print "123", 123;
		} else if (false) {
			print "123", 123;
		}
	*/
	toks := []*token.Token{
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.TRUE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.ELSE},
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.FALSE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingIfStmtWithIfElseElse(t *testing.T) {
	/*
		if (true) {
			print "123", 123;
		} else if (false) {
			print "123", 123;
		} else {
			print "123", 123;
		}
	*/
	toks := []*token.Token{
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.TRUE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.ELSE},
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.FALSE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.ELSE},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingIfStmtWithIfElseIfElseElse(t *testing.T) {
	/*
		if (true) {
			print "123", 123;
		} else if (false) {
			print "123", 123;
		} else if (false) {
			print "123", 123;
		} else {
			print "123", 123;
		}
	*/
	toks := []*token.Token{
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.TRUE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.ELSE},
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.FALSE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.ELSE},
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.FALSE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.ELSE},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingIfStmtWithElseIfElseFail(t *testing.T) {
	/*
		if (true) {
			print "123", 123;
		} else {
			print "123", 123;
		} else if (false) {
			print "123", 123;
		}
	*/
	toks := []*token.Token{
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.TRUE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.ELSE},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.ELSE},
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.FALSE},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.PRINT},
		&token.Token{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		&token.Token{TokenType: token.COMMA},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "123"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParsingForStmtWithAllEmpty(t *testing.T) {
	// for (;;) {}
	toks := []*token.Token{
		&token.Token{TokenType: token.FOR},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithAssignmentInit(t *testing.T) {
	// for (i = 0; i < 10; i = i + 1) {}
	toks := []*token.Token{
		&token.Token{TokenType: token.FOR},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "0"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.LESS},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "10"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.ADD},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "1"},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithDeclarationInit(t *testing.T) {
	// for (let i: int = 0; i < 10; i = i + 1) {}
	toks := []*token.Token{
		&token.Token{TokenType: token.FOR},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "0"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.LESS},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "10"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.ADD},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "1"},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithInitEmpty(t *testing.T) {
	// for (; i < 10; ) {}
	toks := []*token.Token{
		&token.Token{TokenType: token.FOR},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.LESS},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "10"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithConditionEmpty(t *testing.T) {
	// for (let i: int = 0;;i = i + 1) {}
	toks := []*token.Token{
		&token.Token{TokenType: token.FOR},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "0"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.ADD},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "1"},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithIterationEmpty(t *testing.T) {
	// for (let i: int = 0; i < 10;) {}
	toks := []*token.Token{
		&token.Token{TokenType: token.FOR},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "0"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.LESS},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "10"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingWhileStmt(t *testing.T) {
	// while (i < 3) {}
	toks := []*token.Token{
		&token.Token{TokenType: token.WHILE},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.LESS},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "3"},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingWhileStmtWithInvalidConditionExpr(t *testing.T) {
	// while (let i: int = 0) {}
	toks := []*token.Token{
		&token.Token{TokenType: token.WHILE},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "i"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "0"},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParsingWhileStmtWithEmptyConditionExpr(t *testing.T) {
	// while () {}
	toks := []*token.Token{
		&token.Token{TokenType: token.WHILE},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestBreakNode(t *testing.T) {
	// break;
	toks := []*token.Token{
		&token.Token{TokenType: token.BREAK},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestBreakNodeInvalid(t *testing.T) {
	// break let;
	toks := []*token.Token{
		&token.Token{TokenType: token.BREAK},
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestSwitchStmt1(t *testing.T) {
	/*
		switch (a) {
		case 0:
			b = 5;
			break;
		case 1:
			b = 6;
			break;
		case 2:
		case 3:
			b = 7;
		case 4:
			break;
		case 5:
			if (b > 8) {
				b = 10;
			}
		case 1+1:
			b = 7;
			break;
		default:
			b = 9;
		}
	*/
	toks := []*token.Token{
		&token.Token{TokenType: token.SWITCH},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "a"},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.CASE},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "0"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "b"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "5"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.BREAK},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.CASE},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "1"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "b"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "6"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.BREAK},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.CASE},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "2"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.CASE},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "3"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "b"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "7"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.BREAK},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.CASE},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "4"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.BREAK},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.CASE},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "5"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.IF},
		&token.Token{TokenType: token.LEFT_PAREN},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "b"},
		&token.Token{TokenType: token.GREATER},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "8"},
		&token.Token{TokenType: token.RIGHT_PAREN},
		&token.Token{TokenType: token.LEFT_CURLY_BRACE},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "b"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "10"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.CASE},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "1"},
		&token.Token{TokenType: token.ADD},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "1"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "b"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "7"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.BREAK},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.DEFAULT},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "b"},
		&token.Token{TokenType: token.ASSIGN},
		&token.Token{TokenType: token.INT_LITERAL, Raw: "9"},
		&token.Token{TokenType: token.SEMI},
		&token.Token{TokenType: token.RIGHT_CURLY_BRACE},
		&token.Token{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}
