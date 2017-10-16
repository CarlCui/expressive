package parser

import (
	"fmt"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/logger"
	"github.com/carlcui/expressive/scanner"
	"github.com/carlcui/expressive/signature"
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

	parser.expect(token.SEMI)

	return &node
}

func (parser *Parser) isAssignmentStmtStart(tok *token.Token) bool {
	return tok.TokenType == token.IDENTIFIER
}

func (parser *Parser) parseAssignmentStmt() ast.Node {
	if !parser.isAssignmentStmtStart(parser.cur) {
		return parser.syntaxErrorNode("assignment statement")
	}

	var node ast.AssignmentNode
	node.BaseNode = ast.CreateBaseNode(parser.cur, nil)

	identifier := parser.parseIdentifier()
	identifier.SetParent(&node)

	node.Identifier = identifier

	parser.expect(token.ASSIGN)

	expr := parser.parseExpr()
	expr.SetParent(&node)

	node.Expr = expr

	parser.expect(token.SEMI)

	return &node
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

	parser.expect(token.PRINT)

	expr := parser.parseExpr()
	expr.SetParent(&node)

	node.Expr = expr

	parser.expect(token.SEMI)

	return &node
}

// Exprs

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

func (parser *Parser) isTypeLiteralStart(tok *token.Token) bool {
	return tok.TokenType == token.INT_KEYWORD ||
		tok.TokenType == token.FLOAT_KEYWORD ||
		tok.TokenType == token.CHAR_KEYWORD ||
		tok.TokenType == token.STRING_KEYWORD ||
		tok.TokenType == token.BOOL_KEYWORD
}

func (parser *Parser) parseTypeLiteral() ast.Node {
	if !parser.isTypeLiteralStart(parser.cur) {
		return parser.syntaxErrorNode("type literal")
	}

	var node ast.TypeLiteralNode
	node.BaseNode = ast.CreateBaseNode(parser.cur, nil)

	parser.read()

	return &node
}

func (parser *Parser) parseLiteral() ast.Node {
	cur := parser.cur

	if parser.isIntegerLiteralStart(cur) {
		return parser.parserInt()
	} else if parser.isFLoatLiteralStart(cur) {
		return parser.parseFloat()
	} else if parser.isIdentifierStart(cur) {
		return parser.parseIdentifier()
	} else if parser.isBooleanLiteralStart(cur) {
		return parser.parseBool()
	} else if parser.isStringLiteralStart(cur) {
		return parser.parseString()
	} else if parser.isCharacterLiteralStart(cur) {
		return parser.parseCharacter()
	}

	return parser.syntaxErrorNode("literal")
}

func (parser *Parser) parserInt() ast.Node {
	if parser.cur.TokenType != token.INT_LITERAL {
		return parser.syntaxErrorNode("int")
	}

	var node ast.IntegerNode
	node.Init(parser.cur)

	parser.read()

	return &node
}

func (parser *Parser) parseFloat() ast.Node {
	if parser.cur.TokenType != token.FLOAT_LITERAL {
		return parser.syntaxErrorNode("float")
	}

	var node ast.FloatNode
	node.Init(parser.cur)

	parser.read()

	return &node
}

func (parser *Parser) parseString() ast.Node {
	if parser.cur.TokenType != token.STRING_LITERAL {
		return parser.syntaxErrorNode("string")
	}

	var node ast.StringNode
	node.Init(parser.cur)

	parser.read()

	return &node
}

func (parser *Parser) parseCharacter() ast.Node {
	if parser.cur.TokenType != token.CHAR_LITERAL {
		return parser.syntaxErrorNode("character")
	}

	var node ast.CharacterNode
	node.Init(parser.cur)

	parser.read()

	return &node
}

func (parser *Parser) parseBool() ast.Node {
	if !parser.isBooleanLiteralStart(parser.cur) {
		return parser.syntaxErrorNode("boolean")
	}

	var node ast.BooleanNode
	node.Init(parser.cur)

	parser.read()

	return &node
}

func (parser *Parser) parseIdentifier() ast.Node {
	if parser.cur.TokenType != token.IDENTIFIER {
		return parser.syntaxErrorNode("identifier")
	}

	node := ast.IdentifierNode{BaseNode: ast.CreateBaseNode(parser.cur, nil)}

	parser.read()

	return &node
}

func (parser *Parser) syntaxErrorNode(expected string) ast.Node {
	var node ast.ErrorNode
	node.BaseNode = ast.CreateBaseNode(parser.cur, nil)

	node.Expected = expected

	parser.logger.Log(parser.cur.GetLocation(), "expected "+expected)

	return &node
}

func (parser *Parser) isLiteralStart(tok *token.Token) bool {
	return parser.isIntegerLiteralStart(tok) ||
		parser.isFLoatLiteralStart(tok) ||
		parser.isStringLiteralStart(tok) ||
		parser.isCharacterLiteralStart(tok) ||
		parser.isIdentifierStart(tok) ||
		parser.isBooleanLiteralStart(tok)
}

func (parser *Parser) isIntegerLiteralStart(tok *token.Token) bool {
	return parser.cur.TokenType == token.INT_LITERAL
}

func (parser *Parser) isFLoatLiteralStart(tok *token.Token) bool {
	return parser.cur.TokenType == token.FLOAT_LITERAL
}

func (parser *Parser) isStringLiteralStart(tok *token.Token) bool {
	return parser.cur.TokenType == token.STRING_LITERAL
}

func (parser *Parser) isCharacterLiteralStart(tok *token.Token) bool {
	return parser.cur.TokenType == token.CHAR_LITERAL
}

func (parser *Parser) isIdentifierStart(tok *token.Token) bool {
	return parser.cur.TokenType == token.IDENTIFIER
}

func (parser *Parser) isBooleanLiteralStart(tok *token.Token) bool {
	return parser.cur.TokenType == token.TRUE || parser.cur.TokenType == token.FALSE
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
