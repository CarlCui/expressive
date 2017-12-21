package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// IfStmtNode represents a node with if statement
type IfStmtNode struct {
	*BaseNode
	ConditionExprs  []Node
	ConditionBlocks []Node
	ElseBlock       Node
}

// Accept is part of visitor pattern.
func (node *IfStmtNode) Accept(visitor Visitor) {
	visitor.VisitEnterIfStmtNode(node)
	node.VisitChildren(visitor)
	visitor.VisitLeaveIfStmtNode(node)
}

// VisitChildren is part of visitor pattern. Visit left-hand side node, then right-hand side node.
func (node *IfStmtNode) VisitChildren(visitor Visitor) {
	for i, expr := range node.ConditionExprs {
		block := node.ConditionBlocks[i]

		Accept(expr, visitor)
		Accept(block, visitor)
	}

	if node.ElseBlock != nil {
		Accept(node.ElseBlock, visitor)
	}
}

func (node *IfStmtNode) AddCondition(expr Node, block Node) {
	expr.SetParent(node)
	block.SetParent(node)

	node.ConditionExprs = append(node.ConditionExprs, expr)
	node.ConditionBlocks = append(node.ConditionBlocks, block)
}

func (node *IfStmtNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType        string
		Token           *token.Token
		Typing          typing.Typing
		ConditionExprs  []Node
		ConditionBlocks []Node
		ElseBlock       Node
	}{
		NodeType:        "if statement",
		Token:           node.BaseNode.Tok,
		Typing:          node.Typing,
		ConditionExprs:  node.ConditionExprs,
		ConditionBlocks: node.ConditionBlocks,
		ElseBlock:       node.ElseBlock,
	})
}

func CreateIfStmtNode(tok *token.Token) *IfStmtNode {
	var node IfStmtNode
	node.BaseNode = CreateBaseNode(tok, nil)

	node.ConditionExprs = make([]Node, 0)
	node.ConditionBlocks = make([]Node, 0)

	return &node
}
