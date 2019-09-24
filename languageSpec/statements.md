# Statements

## Productions

_stmt_ := _variableDeclarationStmt_ | _assignmentStmt_ | _printStmt_ | _ifStmt_ | _forStmt_ | _whileStmt_ | _switchStmt_ | _breakStmt_

_variableDeclarationStmt_ := (`let`|`const`) _identifier_ _typeAnnotation_? (`=` _expr_)? `;`

_assignmentStmt_ := _expr_ `=` _expr_ `;`

_printStmt_ := `print` _expr_ (`,` _expr_)* `;`
