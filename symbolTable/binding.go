package symbolTable

import (
	"github.com/carlcui/expressive/locator"
	"github.com/carlcui/expressive/typing"
)

// Binding reperents the location in memory of an identifier
type Binding struct {
	IsVariable    bool
	CanBeShadowed bool
	locator       locator.Locator
	typing        typing.Typing
}

func CreateBinding(locator locator.Locator, typing typing.Typing) *Binding {
	return &Binding{true, true, locator, typing}
}

func CreateBindingCannotBeShadowed(locator locator.Locator, typing typing.Typing) *Binding {
	return &Binding{true, false, locator, typing}
}

func (binding *Binding) GetTyping() typing.Typing {
	return binding.typing
}
