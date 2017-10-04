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
	toks := []*token.Token{
		&token.Token{TokenType: token.LET},
		&token.Token{TokenType: token.IDENTIFIER, Raw: "bar"},
		&token.Token{TokenType: token.COLON},
		&token.Token{TokenType: token.INT_KEYWORD},
		&token.Token{TokenType: token.SEMI},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}
