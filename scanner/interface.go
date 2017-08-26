package scanner

import "github.com/carlcui/expressive/token"

// Scanner interface
type Scanner interface {
	Next() *token.Token
}
