package typing

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type PrimitiveType int

const (
	INT PrimitiveType = iota
	FLOAT
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
	CHAR:       "CHAR",
	STRING:     "STRING",
	BOOL:       "BOOL",
	VOID:       "VOID",
	ERROR_TYPE: "ERROR",
	NO_TYPE:    "",
}

var irTypes = [...]string{
	INT:        "i32",
	FLOAT:      "double",
	CHAR:       "i8*",
	STRING:     "i8*",
	BOOL:       "i1",
	VOID:       "void",
	ERROR_TYPE: "",
	NO_TYPE:    "",
}

var sizes = [...]int{
	INT:        4,
	FLOAT:      8,
	CHAR:       4,
	STRING:     4,
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

func (primitiveType PrimitiveType) IrType() string {
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
