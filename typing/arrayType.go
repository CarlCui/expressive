package typing

import "github.com/llir/llvm/ir/types"

// ArrayType implements a basic data type array.
// The sub-type of an array is the type of its elements.
type ArrayType struct {
	SubType Typing
	irType  types.Type
}

// Equals is to compare if this array type is the same as typing.
// Two arrays are considered equal, if their sub-types are the same.
func (arrayType *ArrayType) Equals(typing Typing) bool {
	if a2, ok := typing.(*ArrayType); ok {
		return arrayType.SubType.Equals(a2.SubType)
	}

	return false
}

func (arrayType *ArrayType) String() string {
	return arrayType.SubType.String() + "[]"
}

func (arrayType *ArrayType) Size() int {
	return 4
}

func (arrayType *ArrayType) IrType() types.Type {
	return types.NewPointer(arrayType.SubType.IrType())
}

// func NewArrayType(subType Typing) *ArrayType {
// 	stru := types.NewStruct(types.I32, types.NewArray(10, subType.IrType()))

// }
