package parser

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/scanner"
	"github.com/carlcui/expressive/token"
)

func initParserWithMockTokens(toks []*token.Token) *Parser {
	var scanner scanner.MockScanner

	scanner.Init(toks)

	var parser Parser
	parser.Init(&scanner)

	parser.read()

	return &parser
}

func parseWithMockTokens(toks []*token.Token) ast.Node {
	parser := initParserWithMockTokens(toks)

	return parser.Parse()
}
