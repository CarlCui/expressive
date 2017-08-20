package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
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

	root := parser.Parse()

	fmt.Println(reflect.TypeOf(root))

	b, err := json.MarshalIndent(root, "", "    ")

	if err != nil {
		fmt.Println("error: ", err)
	}

	os.Stdout.Write(b)
	fmt.Println()
}
