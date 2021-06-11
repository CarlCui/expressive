package parser

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/token"
)

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

	var typeNode ast.TypeLiteralNode
	typeNode.BaseNode = ast.CreateBaseNode(parser.cur, nil)

	parser.read()

	if !parser.isArrayTypeStart(parser.cur) {
		return &typeNode
	}

	return parser.parseArrayType(&typeNode)
}

func (parser *Parser) isArrayTypeStart(tok *token.Token) bool {
	return tok.TokenType == token.LEFT_SQUARE_BRACKET
}

func (parser *Parser) parseArrayType(subType ast.Node) ast.Node {
	if !parser.isArrayTypeStart((parser.cur)) {
		return parser.syntaxErrorNode("array type")
	}

	node := ast.CreateArrayTypeNode(parser.cur, subType)

	parser.read()

	parser.expect(token.RIGHT_SQUARE_BRACKET)

	if !parser.isArrayTypeStart(parser.cur) {
		return node
	}

	return parser.parseArrayType(node)
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
