package parser

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/signature"
	"github.com/carlcui/expressive/token"
)

func (parser *Parser) isStmtStart(tok *token.Token) bool {
	return parser.isStmtWithSemiStart(tok) || parser.isStmtWithoutSemiStart(tok)
}

func (parser *Parser) parseStmt() ast.Node {
	if !parser.isStmtStart(parser.cur) {
		return parser.syntaxErrorNode("statement")
	}

	var node ast.Node

	if parser.isStmtWithSemiStart(parser.cur) {
		node = parser.parseStmtWithSemi()
		parser.expect(token.SEMI)
	} else if parser.isStmtWithoutSemiStart(parser.cur) {
		node = parser.parseStmtWithoutSemi()
	}

	if node == nil {
		panic("parseStmt: not a stmt")
	}

	return node
}

func (parser *Parser) isStmtWithSemiStart(tok *token.Token) bool {
	return parser.isVariableDeclarationStmtStart(tok) ||
		parser.isPrintStmtStart(tok) ||
		parser.isBreakStmtStart(tok) ||
		parser.isStmtStartWithExprStart(tok)
}

func (parser *Parser) parseStmtWithSemi() ast.Node {
	if !parser.isStmtWithSemiStart(parser.cur) {
		return parser.syntaxErrorNode("statement with semi")
	}

	var node ast.Node

	if parser.isVariableDeclarationStmtStart(parser.cur) {
		node = parser.parseVariableDeclarationStmt()
	} else if parser.isPrintStmtStart(parser.cur) {
		node = parser.parsePrintStmt()
	} else if parser.isBreakStmtStart(parser.cur) {
		node = parser.parseBreakStmt()
	} else if parser.isStmtStartWithExprStart(parser.cur) {
		node = parser.parseStmtsStartWithExpr()
	}

	return node
}

func (parser *Parser) isStmtWithoutSemiStart(tok *token.Token) bool {
	return parser.isIfStmtStart(tok) ||
		parser.isForStmtStart(tok) ||
		parser.isWhileStmtStart(tok) ||
		parser.isSwitchStmtStart(tok)
}

func (parser *Parser) parseStmtWithoutSemi() ast.Node {
	if !parser.isStmtWithoutSemiStart(parser.cur) {
		return parser.syntaxErrorNode("statement without semi")
	}

	var node ast.Node

	if parser.isIfStmtStart(parser.cur) {
		node = parser.parseIfStmt()
	} else if parser.isForStmtStart(parser.cur) {
		node = parser.parseForStmt()
	} else if parser.isWhileStmtStart(parser.cur) {
		node = parser.parseWhileStmt()
	} else if parser.isSwitchStmtStart(parser.cur) {
		node = parser.parseSwitchStmt()
	}

	return node
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

	identifier.SetParent(&node)

	node.Identifier = identifier

	if parser.cur.TokenType == token.COLON {
		parser.read()

		declaredType := parser.parseTypeLiteral()

		declaredType.SetParent(&node)

		node.DeclaredType = declaredType
	}

	if parser.cur.TokenType == token.ASSIGN {
		parser.read()

		expr := parser.parseExpr()

		expr.SetParent(&node)

		node.Expr = expr
	}

	return &node
}

func (parser *Parser) isStmtStartWithExprStart(tok *token.Token) bool {
	return parser.isExprStart(tok)
}

func (parser *Parser) parseStmtsStartWithExpr() ast.Node {
	parseErrorMsg := "statement with expression start (assignment statement, increment/decrement statement, function call)"

	if !parser.isStmtStartWithExprStart(parser.cur) {
		return parser.syntaxErrorNode(parseErrorMsg)
	}

	leftMostToken := parser.cur

	lhs := parser.parseExpr()

	switch {
	case parser.isAssignmentOperators(parser.cur):
		return parser.parseAssignmentStmt(leftMostToken, lhs)
	case parser.isIncrementDecrementOperator(parser.cur):
		return parser.parseIncrementDecrementStmt(leftMostToken, lhs)
	default:
		return parser.syntaxErrorNode(parseErrorMsg)
	}

}

func (parser *Parser) isAssignmentOperators(tok *token.Token) bool {
	tokenType := tok.TokenType

	return tokenType == token.ASSIGN ||
		tokenType == token.ASSIGN_ADD ||
		tokenType == token.ASSIGN_SUB ||
		tokenType == token.ASSIGN_MUL ||
		tokenType == token.ASSIGN_DIV ||
		tokenType == token.ASSIGN_MOD ||
		tokenType == token.ASSIGN_POW
}

func (parser *Parser) parseAssignmentStmt(baseToken *token.Token, lhs ast.Node) ast.Node {
	if !parser.isAssignmentOperators(parser.cur) {
		return parser.syntaxErrorNode("assignment statement")
	}

	node := ast.CreateAssignmentStmtNode(parser.cur)

	lhs.SetParent(node)
	node.LHS = lhs

	var rhs ast.Node

	switch parser.cur.TokenType {
	case token.ASSIGN:
		parser.expect(token.ASSIGN)

		node.Operator = signature.VOID_OPERATOR

	case token.ASSIGN_ADD, token.ASSIGN_SUB, token.ASSIGN_MUL, token.ASSIGN_DIV, token.ASSIGN_MOD, token.ASSIGN_POW:
		node.Operator = signature.GetOperator(parser.cur)

		parser.read()
	default:
		return parser.syntaxErrorNode("assignment operator")
	}

	rhs = parser.parseExpr()
	rhs.SetParent(node)

	node.RHS = rhs
	return node
}

func (parser *Parser) isIncrementDecrementOperator(tok *token.Token) bool {
	return tok.TokenType == token.INCREMENT || tok.TokenType == token.DECREMENT
}

func (parser *Parser) parseIncrementDecrementStmt(baseToken *token.Token, lhs ast.Node) ast.Node {
	if !parser.isIncrementDecrementOperator(parser.cur) {
		return parser.syntaxErrorNode("increment/decrement statement")
	}

	node := ast.CreateIncDecNode(baseToken)

	if parser.cur.TokenType == token.INCREMENT {
		node.IsIncrement = true
	} else {
		node.IsIncrement = false
	}

	parser.read()

	node.LHS = lhs
	lhs.SetParent(node)

	return node
}

func (parser *Parser) isIfStmtStart(tok *token.Token) bool {
	return tok.TokenType == token.IF
}

func (parser *Parser) parseIfStmt() ast.Node {
	if !parser.isIfStmtStart(parser.cur) {
		return parser.syntaxErrorNode("if statement")
	}

	node := ast.CreateIfStmtNode(parser.cur)

	parser.read()

	parser.expect(token.LEFT_PAREN)
	expr := parser.parseExpr()
	parser.expect(token.RIGHT_PAREN)

	block := parser.parseBlockWithBraces()

	node.AddCondition(expr, block)

	for parser.cur.TokenType == token.ELSE {
		parser.read()

		lastElse := false

		if parser.cur.TokenType != token.IF {
			lastElse = true

			block := parser.parseBlockWithBraces()

			node.ElseBlock = block
			block.SetParent(node)

		} else if parser.cur.TokenType == token.IF {
			parser.read()

			parser.expect(token.LEFT_PAREN)
			expr := parser.parseExpr()
			parser.expect(token.RIGHT_PAREN)

			block := parser.parseBlockWithBraces()

			node.AddCondition(expr, block)
		} else {
			parser.expect(token.IF, token.LEFT_CURLY_BRACE)
		}

		if lastElse {
			break
		}
	}

	return node
}

func (parser *Parser) isWhileStmtStart(tok *token.Token) bool {
	return tok.TokenType == token.WHILE
}

func (parser *Parser) parseWhileStmt() ast.Node {
	if !parser.isWhileStmtStart(parser.cur) {
		return parser.syntaxErrorNode("while statement")
	}

	node := ast.CreateWhileStmtNode(parser.cur)

	parser.read()

	parser.expect(token.LEFT_PAREN)

	node.SetConditionExprNode(parser.parseExpr())

	parser.expect(token.RIGHT_PAREN)

	body := parser.parseBlockWithBraces()

	node.SetBlockNode(body)

	return node
}

func (parser *Parser) isForStmtStart(tok *token.Token) bool {
	return tok.TokenType == token.FOR
}

func (parser *Parser) parseForStmt() ast.Node {
	if !parser.isForStmtStart(parser.cur) {
		return parser.syntaxErrorNode("for statement")
	}

	node := ast.CreateForStmtNode(parser.cur)

	parser.read()

	// condition
	parser.expect(token.LEFT_PAREN)

	if parser.isStmtWithSemiStart(parser.cur) {
		node.SetInitializationStmtNode(parser.parseStmtWithSemi())
	}

	parser.expect(token.SEMI)

	if parser.isExprStart(parser.cur) {
		node.SetConditionExprNode(parser.parseExpr())
	}
	parser.expect(token.SEMI)

	if parser.isStmtWithSemiStart(parser.cur) {
		node.SetIterationStmtNode(parser.parseStmtWithSemi())
	}
	parser.expect(token.RIGHT_PAREN)

	// body
	body := parser.parseBlockWithBraces()

	node.SetBlockNode(body)

	return node
}

func (parser *Parser) isBreakStmtStart(tok *token.Token) bool {
	return tok.TokenType == token.BREAK
}

func (parser *Parser) parseBreakStmt() ast.Node {
	if !parser.isBreakStmtStart(parser.cur) {
		return parser.syntaxErrorNode("break statement")
	}

	node := ast.CreateBreakNode(parser.cur)

	parser.read()

	return node
}

func (parser *Parser) isSwitchStmtStart(tok *token.Token) bool {
	return tok.TokenType == token.SWITCH
}

func (parser *Parser) parseSwitchStmt() ast.Node {
	if !parser.isSwitchStmtStart(parser.cur) {
		return parser.syntaxErrorNode("switch statement")
	}

	node := ast.CreateSwitchStmtNode(parser.cur)

	parser.read()

	parser.expect(token.LEFT_PAREN)

	node.SetTestExpr(parser.parseExpr())

	parser.expect(token.RIGHT_PAREN)

	parser.expect(token.LEFT_CURLY_BRACE)

	if parser.cur.TokenType != token.CASE {
		return parser.syntaxErrorNode("at least one case")
	}

	for parser.cur.TokenType == token.CASE {
		parser.read()

		node.AppendCaseExpr(parser.parseExpr())

		parser.expect(token.COLON)

		node.AppendCaseBlock(parser.parseBlock())
	}

	if parser.cur.TokenType == token.DEFAULT {
		parser.read()
		parser.expect(token.COLON)

		node.SetDefaultBlock(parser.parseBlock())
	}

	parser.expect(token.RIGHT_CURLY_BRACE)

	return node
}

func (parser *Parser) isPrintStmtStart(tok *token.Token) bool {
	return tok.TokenType == token.PRINT
}

func (parser *Parser) parsePrintStmt() ast.Node {
	if !parser.isPrintStmtStart(parser.cur) {
		return parser.syntaxErrorNode("print statement")
	}

	var node ast.PrintNode
	node.BaseNode = ast.CreateBaseNode(parser.cur, nil)

	args := make([]ast.Node, 0)

	parser.expect(token.PRINT)

	stringExpr := parser.parseExpr()
	stringExpr.SetParent(&node)

	node.StringExpr = stringExpr

	for parser.cur.TokenType == token.COMMA {
		parser.read()

		arg := parser.parseExpr()
		arg.SetParent(&node)

		args = append(args, arg)
	}

	node.Args = args

	return &node
}
