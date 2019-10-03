package ast

import (
	"encoding/json"
	"unicode/utf8"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// CharacterNode represents a character node.
type CharacterNode struct {
	*BaseNode
	Val rune
}

// Accept is part of visitor pattern.
func (node *CharacterNode) Accept(visitor Visitor) {
	visitor.VisitCharacterNode(node)
}

// Size returns the size of the rune in number of bytes
func (node *CharacterNode) Size() int {
	return utf8.RuneLen(node.Val)
}

// IsASCII checks whether the character node has a value of a ASCII character
func (node *CharacterNode) IsASCII() bool {
	return node.Size() == 1
}

func (node *CharacterNode) StringValue() string {
	return string(node.Val) + "\x00"
}

// Init initializes an char node with a token
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
