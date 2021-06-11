package parser

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/token"
)

func (parser *Parser) parseBlock() ast.Node {

	var node ast.BlockNode
	node.BaseNode = ast.CreateBaseNode(parser.cur, nil)

	stmts := make([]ast.Node, 0)

	for parser.isStmtStart(parser.cur) {
		stmt := parser.parseStmt()

		stmt.SetParent(&node)

		stmts = append(stmts, stmt)
	}

	node.Stmts = stmts

	return &node
}

func (parser *Parser) parseBlockWithBraces() ast.Node {
	parser.expect(token.LEFT_CURLY_BRACE)
	node := parser.parseBlock()
	parser.expect(token.RIGHT_CURLY_BRACE)

	return node
}
