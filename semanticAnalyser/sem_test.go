package semanticAnalyser

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/input"
	"github.com/carlcui/expressive/logger"
	"github.com/carlcui/expressive/parser"
	"github.com/carlcui/expressive/scanner"
)

func parseFile(dirName string, fileName string) ast.Node {
	var fileInput input.File
	fileInput.Init(dirName, fileName)

	var s scanner.ExpressiveScanner
	s.Init(&fileInput)

	var p parser.Parser
	p.Init(&s)

	return p.Parse()
}

func newLogger() *logger.StdError {
	var logger logger.StdError
	return &logger
}

func parseAndAnalyze(dirName string, fileName string, handleResult func(logger logger.Logger)) {
	root := parseFile(dirName, fileName)

	ast.PrintAst(root)

	logger := newLogger()

	Analyze(root, logger)

	handleResult(logger)
}

func TestAnalzingCorrectPrograms(t *testing.T) {
	dirName := "./testFiles/correct"

	files, err := ioutil.ReadDir(dirName)

	if err != nil {
		panic("Incorrect test file directory!")
	}

	for _, file := range files {
		fileName := file.Name()

		parseAndAnalyze(dirName, fileName, func(logger logger.Logger) {
			if logger.ErrorsCount() > 0 {
				t.Errorf("File %v: error(s) encountered: %v", fileName, logger.ErrorsCount())
			} else {
				fmt.Printf("%v: passed\n", fileName)
			}
		})
	}
}

func TestAnalzingIncorrectPrograms(t *testing.T) {
	dirName := "./testFiles/incorrect"

	files, err := ioutil.ReadDir(dirName)

	if err != nil {
		panic("Incorrect test file directory!")
	}

	for _, file := range files {
		fileName := file.Name()

		parseAndAnalyze(dirName, fileName, func(logger logger.Logger) {
			if logger.ErrorsCount() == 0 {
				t.Errorf("File %v: error not found", fileName)
			}
		})
	}
}

func TestParticularFile(t *testing.T) {
	t.Skip("for local debugging only")

	dirName := "./testFiles/incorrect"
	fileName := "operator_7.exp"

	parseAndAnalyze(dirName, fileName, func(logger logger.Logger) {
		if logger.ErrorsCount() == 0 {
			t.Errorf("File %v: error not found", fileName)
		}
	})
}
