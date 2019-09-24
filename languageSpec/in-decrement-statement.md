# In/Decrement Statement

## Introduction

As use cases for increment and decrement operations in expressions are very limited, and could cause confusion in terms of the position of the operator (e.g. `i++` vs `++i`), they are only implemented as a statement with postfix.

## Production

_incrementStmt_ := _expr_ `++` `;`

_decrementStmt_ := _expr_ `--` `;`