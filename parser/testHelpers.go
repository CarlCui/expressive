package parser

import (
	"testing"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/logger"
	"github.com/carlcui/expressive/scanner"
	"github.com/carlcui/expressive/token"
)

func initParserWithMockTokens(toks []*token.Token) *Parser {
	var scanner scanner.MockScanner
	var logger logger.StdError

	scanner.Init(toks)

	var parser Parser
	parser.Init(&scanner, &logger)

	parser.read()

	return &parser
}

func parseWithMockTokens(toks []*token.Token, handler func(logger logger.Logger)) ast.Node {
	var scanner scanner.MockScanner
	var logger logger.StdError

	scanner.Init(toks)

	var parser Parser
	parser.Init(&scanner, &logger)

	root := parser.Parse()

	handler(&logger)

	return root
}

func reportTestError(message string, node ast.Node, t *testing.T) {
	t.Error(message)
	t.Error(string(ast.SerializeAst(node)))
}

func shouldHaveNoError(t *testing.T) func(logger logger.Logger) {
	return func(logger logger.Logger) {
		if logger.ErrorsCount() > 0 {
			t.Error("Expecting no error.")
		}
	}
}

func shouldHaveError(t *testing.T) func(logger logger.Logger) {
	return func(logger logger.Logger) {
		if logger.ErrorsCount() == 0 {
			t.Error("Expecting error(s).")
		}
	}
}
