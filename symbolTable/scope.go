package symbolTable

import (
	"strconv"

	"github.com/carlcui/expressive/locator"
	"github.com/carlcui/expressive/typing"
)

var nextScopeIndex = 0

// Scope is where variables live and can be referenced
type Scope struct {
	scopeIndex  int
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
	scope.scopeIndex = nextScopeIndex

	nextScopeIndex++

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

// CreateBindingCannotBeShadowed creates a binding associated with an identifier that cannot be shadowed in descendent scopes
func (scope *Scope) CreateBindingCannotBeShadowed(identifier string, locator locator.Locator, typing typing.Typing) *Binding {
	binding := CreateBindingCannotBeShadowed(locator, typing)

	scope.symbolTable.Install(identifier, binding)

	return binding
}

func (scope *Scope) FindBinding(identifier string) *Binding {
	return scope.symbolTable.Lookup(identifier)
}

func (scope *Scope) VariableDeclared(identifier string) bool {
	return scope.FindBinding(identifier) != nil
}

func (scope *Scope) VariableCanBeShadowed(identifier string) bool {
	localScope := scope

	for localScope != nil {
		binding := localScope.FindBinding(identifier)

		if binding != nil && binding.CanBeShadowed == false {
			return false
		}

		localScope = localScope.BaseScope
	}

	return true
}

func (scope *Scope) GetScopeIdentifier() string {
	return "___scope___" + strconv.Itoa(scope.scopeIndex)
}
