package semanticAnalyser

import "github.com/carlcui/expressive/ast"

func Analyze(node ast.Node) {
	var visitor SemanticAnalysisVisitor

	node.Accept(&visitor)
}
