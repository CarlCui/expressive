# Operators

## Logic operators

1. _expr_ `&&` _expr_: and
1. _expr_ `||` _expr_: or
1. `!` _expr_: not

## Arithmetic operators

1. _expr_ `+` _expr_: add
1. _expr_ `-` _expr_: subtract
1. _expr_ `*` _expr_: multiply
1. _expr_ `/` _expr_: divide
1. _expr_ `%` _expr_: remainder
1. _expr_ `^^` _expr_: exponential

## Comparison operators

1. _expr_ `>` _expr_: greater
1. _expr_ `>=` _expr_: greater or equal to
1. _expr_ `<` _expr_: less
1. _expr_ `<=` _expr_: less or equal to
1. _expr_ `==` _expr_: equal (value)
1. _expr_ `===` _expr_: equal (structure)
1. _expr_ `!=` _expr_: not equal (value)
1. _expr_ `!==` _expr_: not equal (structure)

## the ternary operator

1. _expr_ `?` _expr_ `:` _expr_: we all know this one

## type casting

1. `typeof` _identifier_: gives the type of the identifier
1. _identifier_ `instanceof` _type_: type reasoning
1. `cast` _identifier_ `to` _type_: type casting

## function

1. _identifier_ `(` _params_ `)`: function invocation


## optional data type (nullable)

1. _expr_ `?` `.`_identifier_

## array

1. _expr_ `[` (_integerConstant_ | _sliceExpr_) `]`: indexing
1. `...` _expr_: spread operator
1. _expr_ `[` _expr_? `:` _expr_? (`:` _expr_?)? `]`: slicing

| Precedence | operator |
| ---------- | ---------|
| 1          | `.`, `[]`   |
| 2          | `!`, `typeof` |
| 3          | `&&`, `||` |
| 4          | `>`, `<`, `>=`, `<=`, `==`, `!=`, `===`, `!==` |
| 4          | `*`, `/`, `%`, `^^` |
| 5          | `+`, `-`|
| 6          | `? :`        |
