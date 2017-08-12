# Syntax

**Comment is excluded**

_program_ := _importStmts_* (_stmt_ | _exportStmt_)*

_importStmts_ := _importStmt_*

_importStmt_ := `import` _fileName_ (`as` _identifier)? `;`

_fileName_ := _stringExpr_

_exportStmt_ := `export` (_declarationStmt_ | _identifier)

_stmt_ := _declarationStmt_ | _assignmentStmt_

_declarationStmt_ := (`let` | `const`) _identifier

## expressions

_expr_ :=

_exprOr_ := _exprAnd_ (`||` _exprAnd_)*

_exprAnd_ := _exprCmp_ (`&&` _exprCmp_)*

_exprCmp_ := _exprAdd_ ((`>`|`<`|`>=`|`<=`|`==`|`!=`) _exprAdd_)*

_exprAdd_ := _exprMul_ ((`+`|`-`) _exprMul_)*

_exprMul_ := _exprLiteral_ ((`*`|`/`) _exprLiteral_)*

_exprLiteral_ := _constInt_ | _constFloat_ | _identifier_