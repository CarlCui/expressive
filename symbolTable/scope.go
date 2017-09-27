package symbolTable

import "github.com/carlcui/expressive/typing"
import "github.com/carlcui/expressive/locator"

// Scope is where variables live and can be referenced
type Scope struct {
	BaseScope   *Scope
	symbolTable *SymbolTable
}

// CreateScope with a baseScope. If nil, it will use itself as the base scope
func CreateScope(baseScope *Scope) *Scope {
	var scope Scope
	var symbolTable SymbolTable

	symbolTable = make(SymbolTable)

	scope.symbolTable = &symbolTable
	scope.BaseScope = baseScope

	return &scope
}

// CreateSubScope using current scope as base scope
func (scope *Scope) CreateSubScope() *Scope {
	return CreateScope(scope)
}

func (scope *Scope) CreateBinding(identifier string, locator locator.Locator, typing typing.Typing) *Binding {
	binding := CreateBinding(locator, typing)

	scope.symbolTable.Install(identifier, binding)

	return binding
}

func (scope *Scope) FindBinding(identifier string) *Binding {
	return scope.symbolTable.Lookup(identifier)
}

func (scope *Scope) VariableDeclared(identifier string) bool {
	return scope.FindBinding(identifier) != nil
}