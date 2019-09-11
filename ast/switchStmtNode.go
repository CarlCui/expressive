package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// SwitchStmtNode represents a node with for statement
type SwitchStmtNode struct {
	*BaseNode
	TestExpr     Node
	CaseExprs    []Node
	CaseBlocks   []Node
	DefaultBlock Node
	EndLabel     string
}

// Accept is part of visitor pattern.
func (node *SwitchStmtNode) Accept(visitor Visitor) {
	visitor.VisitEnterSwitchStmtNode(node)

	visitor.VisitLeaveSwitchStmtNode(node)
}

func (node *SwitchStmtNode) SetTestExpr(expr Node) {
	node.TestExpr = expr
	expr.SetParent(node)
}

func (node *SwitchStmtNode) AppendCaseExpr(expr Node) {
	node.CaseExprs = append(node.CaseExprs, expr)

	expr.SetParent(node)
}

func (node *SwitchStmtNode) AppendCaseBlock(block Node) {
	node.CaseBlocks = append(node.CaseBlocks, block)
	block.SetParent(node)
}

func (node *SwitchStmtNode) SetDefaultBlock(block Node) {
	node.DefaultBlock = block
	block.SetParent(node)
}

func (node *SwitchStmtNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType     string
		Token        *token.Token
		Typing       typing.Typing
		TestExpr     Node
		CaseExprs    []Node
		CaseBlocks   []Node
		DefaultBlock Node
	}{
		NodeType:     "for statement",
		Token:        node.BaseNode.Tok,
		Typing:       node.Typing,
		TestExpr:     node.TestExpr,
		CaseExprs:    node.CaseExprs,
		CaseBlocks:   node.CaseBlocks,
		DefaultBlock: node.DefaultBlock,
	})
}

func CreateSwitchStmtNode(tok *token.Token) *SwitchStmtNode {
	var node SwitchStmtNode
	node.BaseNode = CreateBaseNode(tok, nil)

	node.CaseExprs = make([]Node, 0)
	node.CaseBlocks = make([]Node, 0)

	return &node
}
