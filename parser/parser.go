package parser

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/scanner"
	"github.com/carlcui/expressive/token"
)

// Parser is a LL1 parser of expressive
type Parser struct {
	scanner *scanner.Scanner

	cur  *token.Token
	prev *token.Token
}

// Init initializes a new parser with given scanner
func (parser *Parser) Init(scanner *scanner.Scanner) {
	parser.scanner = scanner
}

// Parse the current program
func (parser *Parser) Parse() ast.Node {
	return parser.parseProgram()
}

func (parser *Parser) parseProgram() ast.Node {
	return nil
}

// Stmts

func (parser *Parser) parseVariableDeclarationStmt() ast.Node {
	return nil
}

func (parser *Parser) parseAssignmentStmt() ast.Node {
	return nil
}

func (parser *Parser) parsePrintStmt() ast.Node {
	return nil
}

// Exprs

func (parser *Parser) parseExpr() ast.Node {
	return nil
}

func (parser *Parser) parseExprTernaryIfElse() ast.Node {
	return nil
}

func (parser *Parser) parseExprOr() ast.Node {
	return nil
}

func (parser *Parser) parseExprAnd() ast.Node {
	return nil
}

func (parser *Parser) parseExprComp() ast.Node {
	return nil
}

func (parser *Parser) parseExprAdd() ast.Node {
	return nil
}

func (parser *Parser) parseExprMul() ast.Node {
	return nil
}

func (parser *Parser) parseExprNot() ast.Node {
	return nil
}

func (parser *Parser) parseExprFinal() ast.Node {
	return nil
}

func (parser *Parser) parseExprParen() ast.Node {
	return nil
}

func (parser *Parser) parseLiteral() ast.Node {
	if parser.readingInteger() {
		return parser.parserInt()
	} else if parser.readingFloat() {
		return parser.parseFloat()
	} else if parser.readingIdentifier() {
		return parser.parseIdentifier()
	}

	return parser.syntaxErrorNode(token.LITERAL)
}

func (parser *Parser) parserInt() ast.Node {
	if parser.cur.TokenType != token.INT {
		return parser.syntaxErrorNode(token.INT)
	}

	var node ast.IntegerNode
	node.Init(parser.cur)

	parser.read()

	return &node
}

func (parser *Parser) parseFloat() ast.Node {
	if parser.cur.TokenType != token.FLOAT {
		return parser.syntaxErrorNode(token.FLOAT)
	}

	var node ast.FloatNode
	node.Init(parser.cur)

	parser.read()

	return &node
}

func (parser *Parser) parseIdentifier() ast.Node {
	if parser.cur.TokenType != token.IDENTIFIER {
		return parser.syntaxErrorNode(token.IDENTIFIER)
	}

	node := &ast.IdentifierNode{BaseNode: ast.CreateBaseNode(parser.cur, nil)}

	parser.read()

	return node
}

func (parser *Parser) syntaxErrorNode(expected token.Type) ast.Node {
	return nil
}

func (parser *Parser) readingLiteral() bool {
	return parser.readingInteger() || parser.readingFloat() || parser.readingIdentifier()
}

func (parser *Parser) readingInteger() bool {
	return parser.cur.TokenType == token.INT
}

func (parser *Parser) readingFloat() bool {
	return parser.cur.TokenType == token.FLOAT
}

func (parser *Parser) readingIdentifier() bool {
	return parser.cur.TokenType == token.IDENTIFIER
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
