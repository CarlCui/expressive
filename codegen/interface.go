package codegen

import (
	"github.com/carlcui/expressive/ast"
	"github.com/carlcui/expressive/logger"
)

// Generate llvm IR for ast
func Generate(node ast.Node, logger logger.Logger) string {
	var visitor CodegenVisitor
	visitor.Init(logger)

	node.Accept(&visitor)

	rootFragment := visitor.removeVoidFragment(node)

	globalConstants := visitor.constants

	moduleFragment := rootFragment.(*ModuleFragment)

	moduleFragment.Module.Globals = append(globalConstants)

	return moduleFragment.Module.String()
}
