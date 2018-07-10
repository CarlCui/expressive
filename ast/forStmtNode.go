package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// ForStmtNode represents a node with for statement
type ForStmtNode struct {
	*BaseNode
	InitializationStmt Node
	ConditionExpr      Node
	IterationStmt      Node
	Block              Node
}

// Accept is part of visitor pattern.
func (node *ForStmtNode) Accept(visitor Visitor) {
	visitor.VisitEnterForStmtNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveForStmtNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *ForStmtNode) VisitChildren(visitor Visitor) {
	Accept(node.InitializationStmt, visitor)
	Accept(node.ConditionExpr, visitor)
	Accept(node.IterationStmt, visitor)
	Accept(node.Block, visitor)
}

func (node *ForStmtNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType           string
		Token              *token.Token
		Typing             typing.Typing
		InitializationStmt Node
		ConditionExpr      Node
		IterationStmt      Node
		Block              Node
	}{
		NodeType:           "for statement",
		Token:              node.BaseNode.Tok,
		Typing:             node.Typing,
		InitializationStmt: node.InitializationStmt,
		ConditionExpr:      node.ConditionExpr,
		IterationStmt:      node.IterationStmt,
		Block:              node.Block,
	})
}

func CreateForStmtNode(tok *token.Token) *ForStmtNode {
	var node ForStmtNode
	node.BaseNode = CreateBaseNode(tok, nil)

	return &node
}
