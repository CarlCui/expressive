package ast

import (
	"encoding/json"
)

// BlockNode represents a node with assignment statement
type BlockNode struct {
	*BaseNode
	Stmts []Node
}

// Accept is part of visitor pattern.
func (node *BlockNode) Accept(visitor Visitor) {
	visitor.VisitEnterBlockNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveBlockNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *BlockNode) VisitChildren(visitor Visitor) {
	if node.Stmts == nil {
		return
	}

	for _, stmt := range node.Stmts {
		Accept(stmt, visitor)
	}
}

func (node *BlockNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Stmts    []Node
	}{
		NodeType: "block",
		Stmts:    node.Stmts,
	})
}
