package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
	"github.com/llir/llvm/ir"
)

// BreakNode represents a node with break.
type BreakNode struct {
	*BaseNode
	endBlock *ir.Block // the block break should branch to
}

// Accept is part of visitor pattern.
func (node *BreakNode) Accept(visitor Visitor) {
	visitor.VisitBreakNode(node)
}

func (node *BreakNode) FindNearestValidStatementNode() Node {
	ascendentNode := node.GetParent()

	for ascendentNode != nil {
		_, isForStmt := ascendentNode.(*ForStmtNode)
		_, isWhileStmt := ascendentNode.(*WhileStmtNode)
		_, isSwitchStmt := ascendentNode.(*SwitchStmtNode)

		if isForStmt || isWhileStmt || isSwitchStmt {
			node.SetTyping(typing.VOID)
			return ascendentNode
		}

		ascendentNode = ascendentNode.GetParent()
	}

	return nil
}

func (node *BreakNode) FindBreakBlock() *ir.Block {
	nearestValidStatementNode := node.FindNearestValidStatementNode()

	if nearestValidStatementNode == nil {
		panic(node.GetLocation() + "expecting finding valid statement node (for, while or switch)")
	}

	switch stmtNode := nearestValidStatementNode.(type) {
	case *ForStmtNode:
		return stmtNode.EndBlock
	case *WhileStmtNode:
		return stmtNode.EndBlock
	case *SwitchStmtNode:
		return stmtNode.EndBlock
	default:
		panic(node.GetLocation() + "expecting finding valid statement node (for, while or switch)")
	}
}

func (node *BreakNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
	}{
		NodeType: "break",
		Token:    node.BaseNode.Tok,
	})
}

func CreateBreakNode(tok *token.Token) *BreakNode {
	var node BreakNode
	node.BaseNode = CreateBaseNode(tok, nil)

	return &node
}
