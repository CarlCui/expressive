package parser

import (
	"fmt"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/logger"
	"github.com/carlcui/expressive/scanner"
	"github.com/carlcui/expressive/token"
)

// Parser is a LL1 parser of expressive
type Parser struct {
	scanner scanner.Scanner
	logger  logger.Logger

	cur  *token.Token
	prev *token.Token
}

// Init initializes a new parser with given scanner
func (parser *Parser) Init(scanner scanner.Scanner, logger logger.Logger) {
	parser.scanner = scanner
	parser.logger = logger
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

	parser.expect(token.EOF)

	return &node
}

func (parser *Parser) read() {
	parser.prev = parser.cur
	next := parser.scanner.Next()

	for next.TokenType == token.COMMENT {
		next = parser.scanner.Next()
	}

	parser.cur = next
}

func (parser *Parser) expect(tokenTypes ...token.Type) {
	match := false

	for _, tokenType := range tokenTypes {
		if parser.cur.TokenType == tokenType {
			match = true
		}
	}

	if !match {
		expectedMessage := fmt.Sprintf("expected %v", tokenTypes)
		parser.logger.Log(parser.cur.GetLocation(), expectedMessage)
	}

	parser.read()
}
