package symbolTable

// A SymbolTable is a mapping from an identifier to its binding
type SymbolTable map[string]*Binding

func (symbolTable *SymbolTable) Install(identifier string, binding *Binding) {
	if _, ok := (*symbolTable)[identifier]; ok {
		// re-declared, error
	}

	(*symbolTable)[identifier] = binding
}

func (symbolTable *SymbolTable) Lookup(identifier string) *Binding {
	binding, ok := (*symbolTable)[identifier]

	if !ok {
		return nil
	}

	return binding
}
