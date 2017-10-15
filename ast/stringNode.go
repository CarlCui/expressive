package ast

import (
	"encoding/json"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/carlcui/expressive/token"
	"github.com/carlcui/expressive/typing"
)

// StringNode represents an integer constant node.
type StringNode struct {
	*BaseNode
	Val string
}

// Accept is part of visitor pattern.
func (node *StringNode) Accept(visitor Visitor) {
	visitor.VisitStringNode(node)
}

// VisitChildren is part of visitor pattern. Literal node does not have any child.
func (node *StringNode) VisitChildren(visitor Visitor) {

}

// Init initializes an integer node with a token
func (node *StringNode) Init(tok *token.Token) {
	node.BaseNode = CreateBaseNode(tok, nil)

	_, start := utf8.DecodeRuneInString(tok.Raw)
	_, lastSize := utf8.DecodeLastRuneInString(tok.Raw)

	node.Val = tok.Raw[start : len(tok.Raw)-lastSize]
}

func (node *StringNode) EscapeVal() string {
	escapedString := node.Val

	escapedString += "\\00" // append terminating character

	escapedString = strings.Replace(escapedString, "\\n", "\\0A", -1)

	return escapedString
}

func (node *StringNode) EscapedStringLength() int {
	escapedString := node.EscapeVal()

	totalLength := len(escapedString)

	matchEscapedCharacters := regexp.MustCompile("\\\\..")

	escapedCharacters := matchEscapedCharacters.FindAllString(escapedString, -1)

	return totalLength - len(escapedCharacters)*2
}

func (node *StringNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		NodeType string
		Token    *token.Token
		Val      string
		Typing   typing.Typing
	}{
		NodeType: "string literal",
		Token:    node.BaseNode.Tok,
		Val:      node.Val,
		Typing:   node.Typing,
	})
}
