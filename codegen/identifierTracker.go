package codegen

// IdentifierTracker keeps track of identifiers in llvm ir
type IdentifierTracker interface {
	NewIdentifier() string
	CurrentIdentifier() string
	Reset()
}
