package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/symbolTable"
	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// IdentifierNode represents an identifier node.
type IdentifierNode struct {
	*BaseNode
	binding       *symbolTable.Binding
	declaredScope *symbolTable.Scope
}

func (node *IdentifierNode) SetBinding(binding *symbolTable.Binding) {
	node.binding = binding
}

func (node *IdentifierNode) GetBinding() *symbolTable.Binding {
	return node.binding
}

func (node *IdentifierNode) IsBeingDeclared() bool {
	if node.Parent == nil {
		return false
	}

	declarationNode, ok := node.Parent.(*VariableDeclarationNode)

	if !ok {
		return false
	}

	return declarationNode.Identifier == node
}

func (node *IdentifierNode) FindDeclarationScope() *symbolTable.Scope {
	identifier := node.Tok.Raw

	scope := node.GetLocalScope()
	found := scope.VariableDeclared(identifier)

	for !found && scope.BaseScope != nil {
		scope = scope.BaseScope
		found = scope.VariableDeclared(identifier)
	}

	if !found {
		return nil
	}

	return scope
}

func (node *IdentifierNode) FindVariableBinding() *symbolTable.Binding {
	identifier := node.Tok.Raw

	scope := node.FindDeclarationScope()

	if scope == nil {
		return nil
	}

	return scope.FindBinding(identifier)
}

// LocalIdentifier returns the localized identifier name in corresponding scope (appends scope identifier)
func (node *IdentifierNode) LocalIdentifier() string {
	scope := node.FindDeclarationScope()

	if scope == nil {
		panic("cannot get local scope for identifier node")
	}

	return node.BaseNode.Tok.Raw + scope.GetScopeIdentifier()
}

// Accept is part of visitor pattern.
func (node *IdentifierNode) Accept(visitor Visitor) {
	visitor.VisitIdentifierNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *IdentifierNode) VisitChildren(visitor Visitor) {

}

func (node *IdentifierNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Typing   typing.Typing
		Binding  *symbolTable.Binding
	}{
		NodeType: "identifier",
		Token:    node.BaseNode.Tok,
		Typing:   node.Typing,
		Binding:  node.binding,
	})
}
