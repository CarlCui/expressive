package typing

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/llir/llvm/ir/types"
)

type PrimitiveType int

const (
	INT PrimitiveType = iota
	FLOAT
	BYTE
	CHAR
	STRING
	BOOL
	VOID
	ERROR_TYPE
	NO_TYPE
)

var literals = [...]string{
	INT:        "INT",
	FLOAT:      "FLOAT",
	BYTE:       "BYTE",
	CHAR:       "CHAR",
	STRING:     "STRING",
	BOOL:       "BOOL",
	VOID:       "VOID",
	ERROR_TYPE: "ERROR",
	NO_TYPE:    "",
}

var irTypes = [...]types.Type{
	INT:        types.I32,
	FLOAT:      types.Double,
	BYTE:       types.I8,
	CHAR:       types.I8Ptr,
	STRING:     types.I8Ptr,
	BOOL:       types.I1,
	VOID:       types.Void,
	ERROR_TYPE: types.Void,
	NO_TYPE:    types.Void,
}

var sizes = [...]int{
	INT:        4,
	FLOAT:      8,
	BYTE:       1,
	CHAR:       1,
	STRING:     1,
	BOOL:       1,
	VOID:       0,
	ERROR_TYPE: 0,
	NO_TYPE:    0,
}

func (primitiveType PrimitiveType) Equals(typing Typing) bool {
	primitiveType2, ok := typing.(PrimitiveType)

	if !ok {
		return false
	}

	return primitiveType2 == primitiveType
}

func (primitiveType PrimitiveType) Size() int {
	if primitiveType >= 0 && int(primitiveType) < len(literals) {
		return sizes[primitiveType]
	}

	panic(fmt.Sprintf("Illegal primitive type: %v \n", primitiveType))
}

func (primitiveType PrimitiveType) IrType() types.Type {
	if primitiveType >= 0 && int(primitiveType) < len(literals) {
		return irTypes[primitiveType]
	}

	panic(fmt.Sprintf("Illegal primitive type: %v \n", primitiveType))
}

func (primitiveType PrimitiveType) String() string {
	if primitiveType >= 0 && int(primitiveType) < len(literals) {
		return literals[primitiveType]
	}

	return strconv.Itoa(int(primitiveType))
}

func (primitiveType PrimitiveType) MarshalJSON() ([]byte, error) {
	return json.Marshal(primitiveType.String())
}
