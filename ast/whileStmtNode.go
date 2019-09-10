package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// WhileStmtNode represents a node with while statement
type WhileStmtNode struct {
	*BaseNode
	ConditionExpr Node
	Block         Node
}

// Accept is part of visitor pattern.
func (node *WhileStmtNode) Accept(visitor Visitor) {
	visitor.VisitEnterWhileStmtNode(node)
	Accept(node.ConditionExpr, visitor)
	visitor.VisitLeaveWhileStmtNode(node)
}

func (node *WhileStmtNode) SetConditionExprNode(expr Node) {
	node.ConditionExpr = expr
	expr.SetParent(node)
}

func (node *WhileStmtNode) SetBlockNode(block Node) {
	node.Block = block
	block.SetParent(node)
}

func (node *WhileStmtNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType      string
		Token         *token.Token
		Typing        typing.Typing
		ConditionExpr Node
		Block         Node
	}{
		NodeType:      "while statement",
		Token:         node.BaseNode.Tok,
		Typing:        node.Typing,
		ConditionExpr: node.ConditionExpr,
		Block:         node.Block,
	})
}

func CreateWhileStmtNode(tok *token.Token) *WhileStmtNode {
	var node WhileStmtNode
	node.BaseNode = CreateBaseNode(tok, nil)

	return &node
}
