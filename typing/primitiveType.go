package typing

import (
	"encoding/json"
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

func (primitiveType PrimitiveType) Equals(typing Typing) bool {
	primitiveType2, ok := typing.(PrimitiveType)

	if !ok {
		return false
	}

	return primitiveType2 == primitiveType
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
