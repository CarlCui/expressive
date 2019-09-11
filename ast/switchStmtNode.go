package ast

import (
	"encoding/json"
	"strconv"

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

	Accept(node.TestExpr, visitor)
	for i, caseExpr := range node.CaseExprs {
		caseBlock := node.CaseBlocks[i]

		Accept(caseExpr, visitor)
		Accept(caseBlock, visitor)
	}
	Accept(node.DefaultBlock, visitor)

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

func (node *SwitchStmtNode) IsEmptyCaseBlockAt(i int) bool {
	cases := len(node.CaseBlocks)

	if i >= cases {
		panic(node.GetLocation() + " is trying to resolve case block " + strconv.Itoa(i) + " out of range of all case blocks " + strconv.Itoa(cases))
	}

	caseBlock, ok := node.CaseBlocks[i].(*BlockNode)

	if !ok {
		panic(node.CaseBlocks[i].GetLocation() + " is not a block stmt node")
	}

	return caseBlock.IsEmptyBlock()
}

func (node *SwitchStmtNode) IsEmptyDefaultBlock() bool {
	if node.DefaultBlock == nil {
		return true
	}

	defaultBlock, ok := node.DefaultBlock.(*BlockNode)

	if !ok {
		panic(node.DefaultBlock.GetLocation() + " is not a block stmt node")
	}

	return defaultBlock.IsEmptyBlock()
}

func (node *SwitchStmtNode) FindTheNextNonEmptyBlockIndexAt(i int) int {
	cases := len(node.CaseBlocks)

	for ; i < cases; i++ {
		if !node.IsEmptyCaseBlockAt(i) {
			break
		}
	}

	return i
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
