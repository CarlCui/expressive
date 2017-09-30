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
	}

	root := parseWithMockTokens(toks)

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
	}

	root := parseWithMockTokens(toks)

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
