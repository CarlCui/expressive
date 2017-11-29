package ast

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// CharacterNode represents an integer constant node.
type CharacterNode struct {
	*BaseNode
	Val rune
}

// Accept is part of visitor pattern.
func (node *CharacterNode) Accept(visitor Visitor) {
	visitor.VisitCharacterNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *CharacterNode) VisitChildren(visitor Visitor) {

}

// Size returns the size of the rune in number of bytes
func (node *CharacterNode) Size() int {
	return utf8.RuneLen(node.Val)
}

// Init initializes an integer node with a token
func (node *CharacterNode) Init(tok *token.Token) {
	node.BaseNode = CreateBaseNode(tok, nil)

	raw := tok.Raw

	_, totalSize := utf8.DecodeRuneInString(raw)

	firstRune, size := utf8.DecodeRuneInString(raw[totalSize:])

	totalSize += size

	if firstRune == '\\' {
		secondRune, _ := utf8.DecodeRuneInString(raw[totalSize:])

		switch secondRune {
		case 'n':
			node.Val = '\n'
			break
		case 't':
			node.Val = '\t'
			break
		case '0':
			node.Val = '\x00'
			break
		case '\\':
			node.Val = '\\'
			break
		case '\'':
			node.Val = '\''
			break
		case '"':
			node.Val = '"'
			break
		default:
			panic(tok.GetLocation() + ": illegal escape sequence")
		}
	} else {
		node.Val = firstRune
	}
}

func (node *CharacterNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Val      rune
		Typing   typing.Typing
	}{
		NodeType: "character literal",
		Token:    node.BaseNode.Tok,
		Val:      node.Val,
		Typing:   node.Typing,
	})
}
