# Types

## Primitive types

A primitive type including:

1. int
1. float
1. bool
1. char
1. string

## Compound types

### Array type

_arrayType_ := (_arrayType_ | _primitiveType_)`[]`

### Function type

_functionType_ := `func` `(` _type_ (`,` _type_)* `)` _type_

## Type checking

Expressive is a strongly typed language, meaning that it does not do implicit type casting.