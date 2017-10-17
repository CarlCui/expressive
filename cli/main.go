package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/carlcui/expressive/codegen"

	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/input"
	"github.com/carlcui/expressive/logger"
	"github.com/carlcui/expressive/parser"
	"github.com/carlcui/expressive/scanner"
	"github.com/carlcui/expressive/semanticAnalyser"
)

const helpMessage = "Current supported options are: \n" +
	"--asm: produce llvm ir code \n" +
	"--dir/-d: directory of source file \n" +
	"--file/-f: source file name \n" +
	"--help: see all command line options \n" +
	"--outDir: output directory \n"

func parseFile(dirName string, filename string, logger logger.Logger) ast.Node {
	var fileInput input.File
	fileInput.Init(dirName, filename)

	var s scanner.ExpressiveScanner
	s.Init(&fileInput)

	var p parser.Parser
	p.Init(&s, logger)

	root := p.Parse()

	semanticAnalyser.Analyze(root, logger)

	return root
}

func checkSourceFileExtension(filename string) {
	if path.Ext(filename) != ".exp" {
		fmt.Printf("File %v does not end with .exp", filename)
	}
}

func createOutDirIfNotExist(outdir string) {
	if _, err := os.Stat(outdir); err != nil {
		os.MkdirAll(outdir, os.ModeDir|os.ModePerm)
	}
}

func main() {
	args := os.Args

	if len(args) == 1 {
		panic(fmt.Errorf("No options given. \n %v", helpMessage))
	}

	var filename string
	dirName := "." // default to current dir
	outDir := "."  // default to current dir

	options := args[1:]

	for i, arg := range options {

		if arg == "--file" || arg == "-f" {
			if i+1 == len(options) {
				panic("No file specified")
			}

			filename = options[i+1]
		}

		if arg == "--dir" || arg == "-d" {
			if i+1 == len(options) {
				panic("No directory specified")
			}

			dirName = options[i+1]
		}

		if arg == "--outDir" {
			if i+1 == len(options) {
				panic("No output directory specified")
			}

			outDir = options[i+1]
		}
	}

	checkSourceFileExtension(filename)

	var frontendLogger logger.StdError

	root := parseFile(dirName, filename, &frontendLogger)

	if frontendLogger.ErrorsCount() > 0 {
		return
	}

	var codegenLogger logger.StdError

	code := codegen.Generate(root, &codegenLogger)

	if codegenLogger.ErrorsCount() > 0 {
		return
	}

	outputFilename := strings.Replace(filename, ".exp", ".s", -1)

	outfile := path.Join(outDir, outputFilename)

	createOutDirIfNotExist(outDir)

	err := ioutil.WriteFile(outfile, []byte(code), os.ModePerm)

	if err != nil {
		fmt.Println(err)
	}
}
