# Function

## productions

_functionDefinitionStmt_ := `func` (`throwable`)? `(` _formalParameterList_? `)` `->` (_functionReturn_)? _functionDefinitionBlock_

_lambdaDefinitionExpr_ := `(` _formalParameterList_? `)` `->` (_functionReturn)? (_expr_ | _functionDefinitionBlock_)

_formalParameterList_ := _formalParam_ (`,` _formalParam_)*

_formalParam_ := _identifier_ `:` _typeLiteral_

_functionReturn_ := _typeList_

_functionDefinitionBlock_ := `{` _Stmts_ `}`

_functionType_ := `(` _typeList_? `)` `->` _typeLiteral_

_typeList_ := _typeLiteral_ (`,` _typeLiteral_)*

## func

Functions can be defined in a traditional way by using the keyword `func`:

```
func doSomething(input: string) -> int {
    //stmts
}
```

## Throwable

A function has to be marked with keyword `throwable` if it may throw an exception. Otherwise, a compiling error occurs. See detail in `error-handling.md`.

## As a variable

Functions can also be defined as a lambda:

`let someFunc = (value: int) -> (value + 3);`

`let someFunc = (value: int) -> {return value + 3;};`

`let someFunc: (int) -> int = (value: int) -> {return value + 3;};`

`doSomething((value: string) -> parseValue(value));`

**Lambda cannot be `throwable` for simplicity**

If the function type is not specified, the return type will be inferred.

## function type

When declaring a function variable, if no definition is given, you have to specify the function type.

`let someFunc: (int, int) -> int;`