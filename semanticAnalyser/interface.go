package semanticAnalyser

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/logger"
)

func Analyze(node ast.Node, logger logger.Logger) {
	var visitor SemanticAnalysisVisitor
	visitor.logger = logger

	node.Accept(&visitor)
}
