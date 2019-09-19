package typing

import "github.com/llir/llvm/ir/types"

// Typing represents a type in expressive
type Typing interface {
	Equals(typing Typing) bool
	String() string
	Size() int
	IrType() types.Type
}
