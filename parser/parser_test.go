package parser

import (
	"testing"

	"github.com/carlcui/expressive/file"
	"github.com/carlcui/expressive/scanner"
)

func initScanner(fileName string) *scanner.Scanner {
	testFileDir := "."
	testFileName := fileName

	var file file.File
	file.Init(testFileDir, testFileName)

	var scanner scanner.Scanner
	scanner.Init(&file)

	return &scanner
}

func TestParser(t *testing.T) {
	var parser Parser
	parser.Init(initScanner("test1.txt"))

	parser.Parse()
}
