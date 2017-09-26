package semanticAnalyser

import (
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

func TestAnalzingCorrectPrograms(t *testing.T) {
	dirName := "./testFiles/correct"

	files, err := ioutil.ReadDir(dirName)

	if err != nil {
		panic("Incorrect test file directory!")
	}

	for _, file := range files {
		root := parseFile(dirName, file.Name())

		logger := newLogger()

		Analyze(root, logger)

		if logger.ErrorsCount() > 0 {
			t.Errorf("File %v: error(s) encountered: %v", file.Name(), logger.ErrorsCount())
		}
	}
}

func TestAnalzingIncorrectPrograms(t *testing.T) {
	dirName := "./testFiles/incorrect"

	files, err := ioutil.ReadDir(dirName)

	if err != nil {
		panic("Incorrect test file directory!")
	}

	for _, file := range files {
		root := parseFile(dirName, file.Name())

		logger := newLogger()

		Analyze(root, logger)

		if logger.ErrorsCount() == 0 {
			t.Errorf("File %v: error not found", file.Name())
		}
	}
}
