package semanticAnalyser

import (
	"testing"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/input"
	"github.com/carlcui/expressive/logger"
	"github.com/carlcui/expressive/parser"
	"github.com/carlcui/expressive/scanner"
)

func parseFile(fileName string) ast.Node {
	var fileInput input.File
	fileInput.Init("./tests", fileName)

	var s scanner.ExpressiveScanner
	s.Init(&fileInput)

	var p parser.Parser
	p.Init(&s)

	return p.Parse()
}

func TestAnalzingCorrectPrograms(t *testing.T) {
	root := parseFile("correct1.exp")

	var logger logger.StdError

	Analyze(root, logger)

	ast.PrintAst(root)
}
