package ast

import (
	"encoding/json"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// BreakNode represents a node with break.
type BreakNode struct {
	*BaseNode
	label string // the label break should branch to
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

func (node *BreakNode) FindBreakLabel() string {
	nearestValidStatementNode := node.FindNearestValidStatementNode()

	if nearestValidStatementNode == nil {
		panic(node.GetLocation() + "expecting finding valid statement node (for, while or switch)")
	}

	switch stmtNode := nearestValidStatementNode.(type) {
	case *ForStmtNode:
		return stmtNode.EndLabel
	case *WhileStmtNode:
		return stmtNode.EndLabel
	case *SwitchStmtNode:
		return stmtNode.EndLabel
	default:
		panic(node.GetLocation() + "expecting finding valid statement node (for, while or switch)")
	}
}

func (node *BreakNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Label    string
	}{
		NodeType: "break",
		Token:    node.BaseNode.Tok,
		Label:    node.label,
	})
}

func CreateBreakNode(tok *token.Token) *BreakNode {
	var node BreakNode
	node.BaseNode = CreateBaseNode(tok, nil)

	return &node
}
