package parser

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/scanner"
	"github.com/carlcui/expressive/token"
)

type Parser struct {
	scanner *scanner.Scanner

	cur  *token.Token
	prev *token.Token
}

func (parser *Parser) Init(scanner *scanner.Scanner) {
	parser.scanner = scanner
}

func (parser *Parser) Parse() *ast.Node {

}

func (parser *Parser) parseProgram() ast.Node {

}

func (parser *Parser) parseExpr() ast.Node {

}

func (parser *Parser) parseBinaryOperatorExpr() ast.Node {

}

func (parser *Parser) parseLiteral() ast.Node {

}

func (parser *Parser) parserInt() ast.Node {
	if parser.cur.TokenType != token.INT {
		return parser.syntaxErrorNode(token.INT)
	}
}

func (parser *Parser) syntaxErrorNode(expected token.Type) ast.Node {

}

func (parser *Parser) read() {
	parser.prev = parser.cur
	parser.cur = parser.scanner.Next()
}

func (parser *Parser) expect(tokenTypes ...token.Type) {
	for _, tokenType := range tokenTypes {
		if parser.cur.TokenType != tokenType {
			// TODO: error node
		}
	}

	parser.read()
}
