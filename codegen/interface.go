package codegen

import "github.com/carlcui/expressive/ast"
import "github.com/carlcui/expressive/logger"

// Generate llvm IR for ast
func Generate(node ast.Node, logger logger.Logger) string {
	var visitor CodegenVisitor
	visitor.Init(logger)

	node.Accept(&visitor)

	rootFragment := visitor.removeVoidCode(node)

	globalConstants := visitor.constants
	externals := visitor.externals

	rootFragment.AppendBefore(externals)
	rootFragment.AppendBefore(globalConstants)

	return rootFragment.String()
}
