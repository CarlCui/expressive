package ast

import "github.com/carlcui/expressive/token"

// BaseNode has properties that all nodes have
type BaseNode struct {
	tok    *token.Token
	parent *Node
}
