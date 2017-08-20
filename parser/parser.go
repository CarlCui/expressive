package parser

import (
	"fmt"

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
	parser.read()
	return parser.parseProgram()
}

func (parser *Parser) parseProgram() ast.Node {
	var node ast.ProgramNode
	node.Init(parser.cur)

	children := make([]ast.Node, 0)

	for parser.isStmtStart(parser.cur) {
		stmt := parser.parseStmt()

		stmt.SetParent(&node)

		children = append(children, stmt)
	}

	node.Chilren = children

	return &node
}

// Stmts

func (parser *Parser) isStmtStart(tok *token.Token) bool {
	return parser.isVariableDeclarationStmtStart(tok) || parser.isAssignmentStmtStart(tok) || parser.isPrintStmtStart(tok)
}

func (parser *Parser) parseStmt() ast.Node {
	if !parser.isStmtStart(parser.cur) {
		return parser.syntaxErrorNode("statement")
	}

	if parser.isVariableDeclarationStmtStart(parser.cur) {
		return parser.parseVariableDeclarationStmt()
	} else if parser.isAssignmentStmtStart(parser.cur) {
		return parser.parseAssignmentStmt()
	} else if parser.isPrintStmtStart(parser.cur) {
		return parser.parsePrintStmt()
	}

	panic("parseStmt: unreachable")
}

func (parser *Parser) isVariableDeclarationStmtStart(tok *token.Token) bool {
	curTokenType := tok.TokenType

	return curTokenType == token.LET || curTokenType == token.CONST
}

func (parser *Parser) parseVariableDeclarationStmt() ast.Node {
	if !parser.isVariableDeclarationStmtStart(parser.cur) {
		return parser.syntaxErrorNode("variable declaration")
	}

	var node ast.VariableDeclarationNode
	node.BaseNode = ast.CreateBaseNode(parser.cur, nil)

	parser.read()

	identifier := parser.parseIdentifier()

	node.Identifier = identifier

	return &node
}

func (parser *Parser) isAssignmentStmtStart(tok *token.Token) bool {
	return tok.TokenType == token.IDENTIFIER
}

func (parser *Parser) parseAssignmentStmt() ast.Node {
	return nil
}

func (parser *Parser) isPrintStmtStart(tok *token.Token) bool {
	return tok.TokenType == token.PRINT
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

	return parser.syntaxErrorNode("literal")
}

func (parser *Parser) parserInt() ast.Node {
	if parser.cur.TokenType != token.INT {
		return parser.syntaxErrorNode("int")
	}

	var node ast.IntegerNode
	node.Init(parser.cur)

	parser.read()

	return &node
}

func (parser *Parser) parseFloat() ast.Node {
	if parser.cur.TokenType != token.FLOAT {
		return parser.syntaxErrorNode("float")
	}

	var node ast.FloatNode
	node.Init(parser.cur)

	parser.read()

	return &node
}

func (parser *Parser) parseIdentifier() ast.Node {
	if parser.cur.TokenType != token.IDENTIFIER {
		return parser.syntaxErrorNode("identifier")
	}

	node := ast.IdentifierNode{BaseNode: ast.CreateBaseNode(parser.cur, nil)}

	fmt.Println(node)

	parser.read()

	return &node
}

func (parser *Parser) syntaxErrorNode(expected string) ast.Node {
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
