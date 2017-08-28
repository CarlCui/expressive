package parser

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/scanner"
	"github.com/carlcui/expressive/token"
)

func parseWithMockTokens(toks []*token.Token) ast.Node {
	var scanner scanner.MockScanner

	scanner.Init(toks)

	var parser Parser
	parser.Init(&scanner)

	return parser.Parse()
}
