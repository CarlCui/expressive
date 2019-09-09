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
	node.VisitForExpr(visitor)
	visitor.VisitEnterForStmtNodeBeforeBlockNode(node)
	node.VisitForBlock(visitor)
	visitor.VisitLeaveForStmtNode(node)
}

// VisitForExpr traverses nodes in the for conditional expression.
func (node *ForStmtNode) VisitForExpr(visitor Visitor) {
	Accept(node.InitializationStmt, visitor)
	Accept(node.ConditionExpr, visitor)
	Accept(node.IterationStmt, visitor)
}

// VisitForBlock traverses the for block.
func (node *ForStmtNode) VisitForBlock(visitor Visitor) {
	Accept(node.Block, visitor)
}

func (node *ForStmtNode) SetInitializationStmtNode(stmt Node) {
	node.InitializationStmt = stmt
	stmt.SetParent(node)
}

func (node *ForStmtNode) SetConditionExprNode(expr Node) {
	node.ConditionExpr = expr
	expr.SetParent(node)
}

func (node *ForStmtNode) SetIterationStmtNode(stmt Node) {
	node.IterationStmt = stmt
	stmt.SetParent(node)
}

func (node *ForStmtNode) SetBlockNode(block Node) {
	node.Block = block
	block.SetParent(node)
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
