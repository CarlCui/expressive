package parser

import (
	"testing"

	"github.com/carlcui/expressive/token"
)

func TestParseValidVariableDeclarationStmt_LetStmtWithDeclaredArrayType(t *testing.T) {
	// let a: int[];
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.LEFT_SQUARE_BRACKET},
		{TokenType: token.RIGHT_SQUARE_BRACKET},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParseValidVariableDeclarationStmt_LetStmtWithDeclaredArrayOfArrayType(t *testing.T) {
	// let a: int[][];
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.LEFT_SQUARE_BRACKET},
		{TokenType: token.RIGHT_SQUARE_BRACKET},
		{TokenType: token.LEFT_SQUARE_BRACKET},
		{TokenType: token.RIGHT_SQUARE_BRACKET},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveNoError(t))
}

func TestParseInvalidVariableDeclarationStmt_LetStmtWithDeclaredArrayType(t *testing.T) {
	// let a: int[;
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.LEFT_SQUARE_BRACKET},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParseInvalidVariableDeclarationStmt_LetStmtWithDeclaredArrayType2(t *testing.T) {
	// let a: int[][;
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.LEFT_SQUARE_BRACKET},
		{TokenType: token.RIGHT_SQUARE_BRACKET},
		{TokenType: token.LEFT_SQUARE_BRACKET},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParseInvalidVariableDeclarationStmt_LetStmtWithDeclaredArrayType3(t *testing.T) {
	// let a: int[]1;
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.LEFT_SQUARE_BRACKET},
		{TokenType: token.RIGHT_SQUARE_BRACKET},
		{TokenType: token.INT_LITERAL, Raw: "1"},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}

func TestParseInvalidVariableDeclarationStmt_LetStmtWithDeclaredArrayType4(t *testing.T) {
	// let a: int[1];
	toks := []*token.Token{
		{TokenType: token.LET},
		{TokenType: token.IDENTIFIER, Raw: "a"},
		{TokenType: token.COLON},
		{TokenType: token.INT_KEYWORD},
		{TokenType: token.LEFT_SQUARE_BRACKET},
		{TokenType: token.INT_LITERAL, Raw: "1"},
		{TokenType: token.RIGHT_SQUARE_BRACKET},
		{TokenType: token.SEMI},
		{TokenType: token.EOF},
	}

	parseWithMockTokens(toks, shouldHaveError(t))
}
