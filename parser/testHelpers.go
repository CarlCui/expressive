package parser

import (
	"testing"

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

func reportTestError(message string, node ast.Node, t *testing.T) {
	t.Error(message)
	t.Error(string(ast.SerializeAst(node)))
}
