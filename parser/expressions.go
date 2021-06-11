package parser

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/token"
)

func (parser *Parser) isExprStart(tok *token.Token) bool {
	return parser.isExprTernaryIfElseStart(tok)
}

func (parser *Parser) parseExpr() ast.Node {
	if !parser.isExprStart(parser.cur) {
		return parser.syntaxErrorNode("expression")
	}

	return parser.parseExprTernaryIfElse()
}

func (parser *Parser) isExprTernaryIfElseStart(tok *token.Token) bool {
	return parser.isExprOrStart(tok)
}

func (parser *Parser) parseExprTernaryIfElse() ast.Node {
	if !parser.isExprTernaryIfElseStart(parser.cur) {
		return parser.syntaxErrorNode("ternary if else expression")
	}

	expr1 := parser.parseExprOr()

	if parser.cur.TokenType == token.QUESTION_MARK {
		cur := parser.cur

		parser.read()

		expr2 := parser.parseExprOr()

		parser.expect(token.COLON)

		expr3 := parser.parseExprOr()

		expr1 = ast.CreateTernaryOperatorNode(cur, signature.IF_ELSE, expr1, expr2, expr3)
	}

	return expr1
}

func (parser *Parser) isExprOrStart(tok *token.Token) bool {
	return parser.isExprAndStart(tok)
}

func (parser *Parser) parseExprOr() ast.Node {
	if !parser.isExprOrStart(parser.cur) {
		return parser.syntaxErrorNode("logical or expression")
	}

	lhs := parser.parseExprAnd()

	for parser.cur.TokenType == token.LOR {
		cur := parser.cur

		parser.read()

		rhs := parser.parseExprAnd()

		lhs = ast.CreateBinaryOperatorNode(cur, signature.LOGIC_OR, lhs, rhs)
	}
	return lhs
}

func (parser *Parser) isExprAndStart(tok *token.Token) bool {
	return parser.isExprCompStart(tok)
}

func (parser *Parser) parseExprAnd() ast.Node {
	if !parser.isExprAndStart(parser.cur) {
		return parser.syntaxErrorNode("logical and expression")
	}

	lhs := parser.parseExprComp()

	for parser.cur.TokenType == token.LAND {
		cur := parser.cur

		parser.read()

		rhs := parser.parseExprComp()

		lhs = ast.CreateBinaryOperatorNode(cur, signature.GetOperator(cur), lhs, rhs)
	}

	return lhs
}

func (parser *Parser) isExprCompStart(tok *token.Token) bool {
	return parser.isExprAddStart(tok)
}

func (parser *Parser) isExprCompOperator(tok *token.Token) bool {
	tokenType := tok.TokenType

	return tokenType == token.LESS ||
		tokenType == token.LEQ ||
		tokenType == token.GREATER ||
		tokenType == token.GEQ ||
		tokenType == token.EQUAL ||
		tokenType == token.NOT_EQUAL ||
		tokenType == token.TRIPLE_EQUAL ||
		tokenType == token.TRIPLE_NOT_EQUAL
}

func (parser *Parser) parseExprComp() ast.Node {
	if !parser.isExprCompStart(parser.cur) {
		return parser.syntaxErrorNode("comparison expression")
	}

	lhs := parser.parseExprAdd()

	for parser.isExprCompOperator(parser.cur) {
		cur := parser.cur

		parser.read()

		rhs := parser.parseExprAdd()

		lhs = ast.CreateBinaryOperatorNode(cur, signature.GetOperator(cur), lhs, rhs)
	}

	return lhs
}

func (parser *Parser) isExprAddStart(tok *token.Token) bool {
	return parser.isExprMulStart(tok)
}

func (parser *Parser) isExprAddOperator(tok *token.Token) bool {
	tokenType := tok.TokenType

	return tokenType == token.ADD || tokenType == token.SUB
}

func (parser *Parser) parseExprAdd() ast.Node {
	if !parser.isExprAddStart(parser.cur) {
		return parser.syntaxErrorNode("addition expression")
	}

	lhs := parser.parseExprMul()

	for parser.isExprAddOperator(parser.cur) {
		cur := parser.cur

		parser.read()

		rhs := parser.parseExprMul()

		lhs = ast.CreateBinaryOperatorNode(cur, signature.GetOperator(cur), lhs, rhs)
	}

	return lhs
}

func (parser *Parser) isExprMulStart(tok *token.Token) bool {
	return parser.isExprNotStart(tok)
}

func (parser *Parser) isExprMulOperator(tok *token.Token) bool {
	return tok.TokenType == token.MUL || tok.TokenType == token.DIV || tok.TokenType == token.MOD || tok.TokenType == token.POW
}

func (parser *Parser) parseExprMul() ast.Node {
	if !parser.isExprMulStart(parser.cur) {
		return parser.syntaxErrorNode("multiplication expression")
	}

	lhs := parser.parseExprNot()

	for parser.isExprMulOperator(parser.cur) {
		cur := parser.cur

		parser.read()

		rhs := parser.parseExprNot()

		lhs = ast.CreateBinaryOperatorNode(cur, signature.GetOperator(cur), lhs, rhs)
	}

	return lhs
}

func (parser *Parser) isExprNotStart(tok *token.Token) bool {
	return tok.TokenType == token.LNOT || parser.isExprFinalStart(tok)
}

func (parser *Parser) parseExprNot() ast.Node {
	if !parser.isExprNotStart(parser.cur) {
		return parser.syntaxErrorNode("logic not expression")
	}

	if parser.cur.TokenType == token.LNOT {
		currentToken := parser.cur

		parser.read()

		expr := parser.parseExprNot()

		node := ast.CreateUnaryOperatorNode(currentToken, signature.GetOperator(currentToken), expr)

		return node
	}

	return parser.parseExprFinal()
}

func (parser *Parser) isExprFinalStart(tok *token.Token) bool {
	return parser.isExprParenStart(tok) || parser.isLiteralStart(tok)
}

func (parser *Parser) parseExprFinal() ast.Node {
	if !parser.isExprFinalStart(parser.cur) {
		return parser.syntaxErrorNode("final expression (literal or paren)")
	}

	if parser.isExprParenStart(parser.cur) {
		return parser.parseExprParen()
	}

	return parser.parseLiteral()
}

func (parser *Parser) isExprParenStart(tok *token.Token) bool {
	return tok.TokenType == token.LEFT_PAREN
}

func (parser *Parser) parseExprParen() ast.Node {
	if !parser.isExprParenStart(parser.cur) {
		return parser.syntaxErrorNode("parenthesis expression")
	}

	parser.expect(token.LEFT_PAREN)

	expr := parser.parseExpr()

	parser.expect(token.RIGHT_PAREN)

	return expr
}
