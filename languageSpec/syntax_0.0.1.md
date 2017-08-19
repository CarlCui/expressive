# syntax for version 0.0.1

## Features

In v0.0.1, express will have the following features:

1. Arithmetic and logic operations
1. variable declaration, assignment.
1. const
1. basic data types + string
1. print (built-in for debugging purpose)
1. basic type checking
1. basic type inference

The whole file is wrapped in a main function, since no function support yet.

## productions

_program_ := _stmt_*

### Statements

_stmt_ := _variableDeclarationStmt_ | _assignmentStmt_ | _printStmt_

_variableDeclarationStmt_ := (`let`|`const`) _identifier_ _typeAnnotation_? (`=` _expr_)? `;`

_assignmentStmt_ := _identifier_ `=` _expr_ `;`

_printStmt_ := `print` `(` _expr_ `)` `;`

### Expressions

_expr_ := _exprTernaryIfElse_

_exprTernaryIfElse_ := _exprOr_ `?` _exprOr_ `: `_exprOr_

_exprOr_ := _exprAnd_ (`||` _exprAnd)*

_exprAnd_ := _exprComp_ (`&&` _exprComp_)*

_exprComp_ := _exprAdd_ (`>`|`<`|`>=`|`<=`|`==`|`!=`|`===`|`!==` _exprAdd_)*

_exprAdd_ := _exprMul_ (`+`|`-` _exprMul_)*

_exprMul_ := _exprNot_ (`*`|`/`|`%`|`^^` _exprNot_)*

_exprNot_ := `!` _exprFinal_

_exprFinal_ := _exprParen_ | _literal_

_exprParen_ := `(` _expr_ `)`


### Literals and misc

_literal_ := _intLiteral_ | _floatLiteral_ | _booleanLiteral_ | _charLiteral_ | _stringLiteral_ | _identifier_

_typeAnnotation_ := `:` _typeLiteral_

_typeLiteral_ := `int` | `bool` | `float` | `char` | `string`

## Operator precedence

| Precedence | operator |
| ---------- | ---------|
| 1          | `()`     |
| 2          | `!` |
| 3          | `*`, `/`, `%`, `^^` |
| 4          | `+`, `-`|
| 5          | `>`, `<`, `>=`, `<=`, `==`, `!=`, `===`, `!==` |
| 6 | `&&` |
| 7 | `||` |
| 8          | `? :` |
