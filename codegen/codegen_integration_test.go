package codegen

import (
	"fmt"
	"testing"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/input"
	"github.com/carlcui/expressive/logger"
	"github.com/carlcui/expressive/parser"
	"github.com/carlcui/expressive/scanner"
	"github.com/carlcui/expressive/semanticAnalyser"
)

func parseFile(dirName string, fileName string) ast.Node {
	var fileInput input.File
	fileInput.Init(dirName, fileName)

	var s scanner.ExpressiveScanner
	s.Init(&fileInput)

	var p parser.Parser
	p.Init(&s, newLogger())

	return p.Parse()
}

func newLogger() *logger.StdError {
	var logger logger.StdError
	return &logger
}

func getAnalyzedAst(dirName string, fileName string) ast.Node {
	root := parseFile(dirName, fileName)

	logger := newLogger()

	semanticAnalyser.Analyze(root, logger)

	return root
}

func TestCodegen(t *testing.T) {
	root := getAnalyzedAst("tests", "test1.exp")

	result := Generate(root, newLogger())

	fmt.Println(result)
}
