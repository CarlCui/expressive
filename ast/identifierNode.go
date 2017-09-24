package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/symbolTable"
	"github.com/carlcui/expressive/token"
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

	_, ok := node.Parent.(*VariableDeclarationNode)

	return ok
}

func (node *IdentifierNode) FindVariableBinding() *symbolTable.Binding {
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

	return scope.FindBinding(identifier)
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
	}{
		NodeType: "identifier",
		Token:    node.BaseNode.Tok,
	})
}
