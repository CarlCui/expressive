package parser

import (
	"testing"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/token"
)

func TestParseValidVariableDelarationStmt_LetStmt(t *testing.T) {
	toks := []*token.Token{
		{TokenType: token.LET, Raw: "let"},
		{TokenType: token.IDENTIFIER, Raw: "foo"},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
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
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "foo"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
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
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.INT_LITERAL, Raw: "0"},
		{TokenType: token.INT_LITERAL, Raw: "-1"},
		{TokenType: token.INT_LITERAL, Raw: "-0"},
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
		{TokenType: token.FLOAT_LITERAL, Raw: "123.123"},
		{TokenType: token.FLOAT_LITERAL, Raw: "123.0"},
		{TokenType: token.FLOAT_LITERAL, Raw: "-1.0"},
		{TokenType: token.FLOAT_LITERAL, Raw: "-123.5"},
		{TokenType: token.FLOAT_LITERAL, Raw: "-0.6"},
		{TokenType: token.FLOAT_LITERAL, Raw: "-0.0"},
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
		{TokenType: token.TRUE},
		{TokenType: token.FALSE},
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
		{TokenType: token.ADD},
		{TokenType: token.IDENTIFIER},
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
		{TokenType: token.CHAR_LITERAL, Raw: "'a'"},
		{TokenType: token.CHAR_LITERAL, Raw: "'1'"},
		{TokenType: token.CHAR_LITERAL, Raw: "'\\0'"},
		{TokenType: token.CHAR_LITERAL, Raw: "'\\''"},
		{TokenType: token.CHAR_LITERAL, Raw: "'\\\"'"},
		{TokenType: token.CHAR_LITERAL, Raw: "'\\t'"},
		{TokenType: token.CHAR_LITERAL, Raw: "'\\n'"},
		{TokenType: token.CHAR_LITERAL, Raw: "'\\\\'"},
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
		{TokenType: token.STRING_LITERAL, Raw: "\"abc\""},
		{TokenType: token.STRING_LITERAL, Raw: "\"\""},
		{TokenType: token.STRING_LITERAL, Raw: "\"  \""},
		{TokenType: token.STRING_LITERAL, Raw: "\"\\\"\""},
		{TokenType: token.STRING_LITERAL, Raw: "\"\\'\""},
		{TokenType: token.STRING_LITERAL, Raw: "\"\\n\""},
		{TokenType: token.STRING_LITERAL, Raw: "\"\\t\""},
		{TokenType: token.STRING_LITERAL, Raw: "\"\\0\""},
		{TokenType: token.STRING_LITERAL, Raw: "\"\\\\\""},
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
		{TokenType: token.COMMENT},
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "foo"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestSkipCommentTokenWhenParsing_InBetween(t *testing.T) {
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "foo"},
		{TokenType: token.COMMENT},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.COMMENT},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestSkipCommentTokenWhenParsing_After(t *testing.T) {
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "foo"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.SEMI},
		{TokenType: token.COMMENT},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestSkipCommentTokenWhenParsing_AfterEOF(t *testing.T) {
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "foo"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
		{TokenType: token.COMMENT},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestNoEOFShouldFail(t *testing.T) {
	// let bar: int; (no EOF)
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "bar"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.SEMI},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParsingAssignmentStmt(t *testing.T) {
	// a = 5;
	toks := []*token.Token{
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "5"},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingAssignmentStmtWithLHSExpression(t *testing.T) {
	// a = 5 + 3;
	toks := []*token.Token{
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "5"},
		{TokenType: token.ADD},
		{TokenType: token.INT_LITERAL, Raw: "3"},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingAssignmentStmtWithLHSEmpty(t *testing.T) {
	// a =;
	toks := []*token.Token{
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.ASSIGN},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParsingCompoundAssignment(t *testing.T) {
	// a += 5 + 3;
	toks := []*token.Token{
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.ASSIGN_ADD},
		{TokenType: token.INT_LITERAL, Raw: "5"},
		{TokenType: token.ADD},
		{TokenType: token.INT_LITERAL, Raw: "3"},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingIncrementStatement(t *testing.T) {
	// a++;
	toks := []*token.Token{
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.INCREMENT},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingDecrementStatement(t *testing.T) {
	// a--;
	toks := []*token.Token{
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.DECREMENT},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingIncrementStatementWithExpr(t *testing.T) {
	// (a+1)++; -- although not semantically correct in expressive, the parser still permits this
	toks := []*token.Token{
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.ADD},
		{TokenType: token.INT_LITERAL, Raw: "1"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.INCREMENT},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingIncrementStatementIncorrect(t *testing.T) {
	// a++5;
	toks := []*token.Token{
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.INCREMENT},
		{TokenType: token.INT_LITERAL, Raw: "5"},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParsingPrintStmt(t *testing.T) {
	// print "123";
	toks := []*token.Token{
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingPrintStmtWithOneArg(t *testing.T) {
	// print "123", 123;
	toks := []*token.Token{
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingPrintStmtWithMultipleArgs(t *testing.T) {
	// print "123", 123, 123.5;
	toks := []*token.Token{
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.COMMA},
		{TokenType: token.FLOAT_LITERAL, Raw: "123.5"},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingPrintStmtWithMultipleArgs_Fail(t *testing.T) {
	// print "123", 123,;
	toks := []*token.Token{
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.COMMA},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
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
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.TRUE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
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
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.TRUE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.ELSE},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
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
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.TRUE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.ELSE},
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.FALSE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
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
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.TRUE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.ELSE},
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.FALSE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.ELSE},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
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
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.TRUE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.ELSE},
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.FALSE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.ELSE},
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.FALSE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.ELSE},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
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
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.TRUE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.ELSE},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.ELSE},
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.FALSE},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.PRINT},
		{TokenType: token.STRING_LITERAL, Raw: "\"123\""},
		{TokenType: token.COMMA},
		{TokenType: token.INT_LITERAL, Raw: "123"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParsingForStmtWithAllEmpty(t *testing.T) {
	// for (;;) {}
	toks := []*token.Token{
		{TokenType: token.FOR},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.SEMI},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithAssignmentInit(t *testing.T) {
	// for (i = 0; i < 10; i = i + 1) {}
	toks := []*token.Token{
		{TokenType: token.FOR},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "0"},
		{TokenType: token.SEMI},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.LESS},
		{TokenType: token.INT_LITERAL, Raw: "10"},
		{TokenType: token.SEMI},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.ASSIGN},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.ADD},
		{TokenType: token.INT_LITERAL, Raw: "1"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithDeclarationInit(t *testing.T) {
	// for (let i: int = 0; i < 10; i = i + 1) {}
	toks := []*token.Token{
		{TokenType: token.FOR},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "0"},
		{TokenType: token.SEMI},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.LESS},
		{TokenType: token.INT_LITERAL, Raw: "10"},
		{TokenType: token.SEMI},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.ASSIGN},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.ADD},
		{TokenType: token.INT_LITERAL, Raw: "1"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithInitEmpty(t *testing.T) {
	// for (; i < 10; ) {}
	toks := []*token.Token{
		{TokenType: token.FOR},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.SEMI},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.LESS},
		{TokenType: token.INT_LITERAL, Raw: "10"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithConditionEmpty(t *testing.T) {
	// for (let i: int = 0;;i = i + 1) {}
	toks := []*token.Token{
		{TokenType: token.FOR},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "0"},
		{TokenType: token.SEMI},
		{TokenType: token.SEMI},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.ASSIGN},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.ADD},
		{TokenType: token.INT_LITERAL, Raw: "1"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingForStmtWithIterationEmpty(t *testing.T) {
	// for (let i: int = 0; i < 10;) {}
	toks := []*token.Token{
		{TokenType: token.FOR},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "0"},
		{TokenType: token.SEMI},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.LESS},
		{TokenType: token.INT_LITERAL, Raw: "10"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingWhileStmt(t *testing.T) {
	// while (i < 3) {}
	toks := []*token.Token{
		{TokenType: token.WHILE},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.LESS},
		{TokenType: token.INT_LITERAL, Raw: "3"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParsingWhileStmtWithInvalidConditionExpr(t *testing.T) {
	// while (let i: int = 0) {}
	toks := []*token.Token{
		{TokenType: token.WHILE},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "i"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "0"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParsingWhileStmtWithEmptyConditionExpr(t *testing.T) {
	// while () {}
	toks := []*token.Token{
		{TokenType: token.WHILE},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestBreakNode(t *testing.T) {
	// break;
	toks := []*token.Token{
		{TokenType: token.BREAK},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestBreakNodeInvalid(t *testing.T) {
	// break let;
	toks := []*token.Token{
		{TokenType: token.BREAK},
		{TokenType: token.LET},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
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
		{TokenType: token.SWITCH},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.CASE},
		{TokenType: token.INT_LITERAL, Raw: "0"},
		{TokenType: token.COLON},
		{TokenType: token.IDENTIFIER, Raw: "b"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "5"},
		{TokenType: token.SEMI},
		{TokenType: token.BREAK},
		{TokenType: token.SEMI},
		{TokenType: token.CASE},
		{TokenType: token.INT_LITERAL, Raw: "1"},
		{TokenType: token.COLON},
		{TokenType: token.IDENTIFIER, Raw: "b"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "6"},
		{TokenType: token.SEMI},
		{TokenType: token.BREAK},
		{TokenType: token.SEMI},
		{TokenType: token.CASE},
		{TokenType: token.INT_LITERAL, Raw: "2"},
		{TokenType: token.COLON},
		{TokenType: token.CASE},
		{TokenType: token.INT_LITERAL, Raw: "3"},
		{TokenType: token.COLON},
		{TokenType: token.IDENTIFIER, Raw: "b"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "7"},
		{TokenType: token.SEMI},
		{TokenType: token.BREAK},
		{TokenType: token.SEMI},
		{TokenType: token.CASE},
		{TokenType: token.INT_LITERAL, Raw: "4"},
		{TokenType: token.COLON},
		{TokenType: token.BREAK},
		{TokenType: token.SEMI},
		{TokenType: token.CASE},
		{TokenType: token.INT_LITERAL, Raw: "5"},
		{TokenType: token.COLON},
		{TokenType: token.IF},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.IDENTIFIER, Raw: "b"},
		{TokenType: token.GREATER},
		{TokenType: token.INT_LITERAL, Raw: "8"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.IDENTIFIER, Raw: "b"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "10"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.CASE},
		{TokenType: token.INT_LITERAL, Raw: "1"},
		{TokenType: token.ADD},
		{TokenType: token.INT_LITERAL, Raw: "1"},
		{TokenType: token.COLON},
		{TokenType: token.IDENTIFIER, Raw: "b"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "7"},
		{TokenType: token.SEMI},
		{TokenType: token.BREAK},
		{TokenType: token.SEMI},
		{TokenType: token.DEFAULT},
		{TokenType: token.COLON},
		{TokenType: token.IDENTIFIER, Raw: "b"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "9"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestSwitchStmtWithoutDefault(t *testing.T) {
	/*
		switch (a) {
		case 0:
			b = 5;
			break;
		}
	*/
	toks := []*token.Token{
		{TokenType: token.SWITCH},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.CASE},
		{TokenType: token.INT_LITERAL, Raw: "0"},
		{TokenType: token.COLON},
		{TokenType: token.IDENTIFIER, Raw: "b"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "5"},
		{TokenType: token.SEMI},
		{TokenType: token.BREAK},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestSwitchStmtWithEmptyCase(t *testing.T) {
	/*
		switch (a) {
		default:
			b = 9;
		}
	*/
	toks := []*token.Token{
		{TokenType: token.SWITCH},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.DEFAULT},
		{TokenType: token.COLON},
		{TokenType: token.IDENTIFIER, Raw: "b"},
		{TokenType: token.ASSIGN},
		{TokenType: token.INT_LITERAL, Raw: "9"},
		{TokenType: token.SEMI},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestSwitchStmtWithEmpty(t *testing.T) {
	/*
		switch (a) {
		}
	*/
	toks := []*token.Token{
		{TokenType: token.SWITCH},
		{TokenType: token.LEFT_PAREN},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.RIGHT_PAREN},
		{TokenType: token.LEFT_CURLY_BRACE},
		{TokenType: token.RIGHT_CURLY_BRACE},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}
