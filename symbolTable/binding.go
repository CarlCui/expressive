package symbolTable

import "github.com/carlcui/expressive/locator"
import "github.com/carlcui/expressive/typing"

// Binding reperents the location in memory of an identifier
type Binding struct {
	IsVariable bool
	locator    locator.Locator
	typing     typing.Typing
}

func CreateBinding(locator locator.Locator, typing typing.Typing) *Binding {
	return &Binding{true, locator, typing}
}

func (binding *Binding) GetTyping() typing.Typing {
	return binding.typing
}
